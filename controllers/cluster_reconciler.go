// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"time"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/internal/config"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/clustermeta"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/clusters"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/util"
)

type ClusterReconciler struct {
	clusters.ManagedReconciler

	clustersManager *clusters.Manager
	config          config.Configuration

	clusterID types.UID
	queue     workqueue.RateLimitingInterface
}

func NewClusterReconciler(name string, log logr.Logger, clustersManager *clusters.Manager, config config.Configuration) *ClusterReconciler {
	return &ClusterReconciler{
		ManagedReconciler: clusters.NewManagedReconciler(name, log),

		clustersManager: clustersManager,
		config:          config,
	}
}

func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.GetLogger().WithValues("cluster", req.NamespacedName)

	err := r.setClusterID(ctx)
	if err != nil {
		return ctrl.Result{}, errors.WithStackIf(err)
	}

	log.Info("reconciling")

	cluster := &clusterregistryv1alpha1.Cluster{}
	err = r.GetClient().Get(ctx, req.NamespacedName, cluster)
	if apierrors.IsNotFound(err) {
		removeErr := r.removeRemoteCluster(req.NamespacedName.Name)
		if removeErr != nil && !errors.Is(errors.Cause(removeErr), clusters.ErrClusterNotFound) {
			return ctrl.Result{}, errors.WithStackIf(removeErr)
		}

		// trigger local reconciles to fix possible local conflict status
		if err = r.triggerLocalClusterReconciles(ctx); err != nil {
			return ctrl.Result{}, errors.WithStackIf(err)
		}

		log.Info("reconciled successfully")

		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "could not get object")
	}

	isClusterLocal := cluster.Spec.ClusterID == r.clusterID

	cluster.Status = cluster.Status.Reset()
	if isClusterLocal {
		// try remove existing remote cluster instance, since a local cluster could have been a remote earlier
		removeErr := r.removeRemoteCluster(cluster.Name)
		if removeErr != nil && !errors.Is(errors.Cause(removeErr), clusters.ErrClusterNotFound) {
			return ctrl.Result{}, errors.WithStackIf(removeErr)
		}
		cluster.Status.Type = clusterregistryv1alpha1.ClusterTypeLocal
	} else if cluster.Status.Type == clusterregistryv1alpha1.ClusterTypeLocal {
		if err := r.triggerLocalClusterReconciles(ctx); err != nil {
			return ctrl.Result{}, errors.WithStackIf(err)
		}
	}

	currentConditions := GetCurrentConditions(cluster)

	var reconcileError error
	if _, ok := cluster.GetAnnotations()[clusterregistryv1alpha1.ClusterDisabledAnnotation]; ok { //nolint:nestif
		removeErr := r.removeRemoteCluster(cluster.Name)
		if removeErr != nil && !errors.Is(errors.Cause(removeErr), clusters.ErrClusterNotFound) {
			return ctrl.Result{}, errors.WithStackIf(removeErr)
		}

		cluster.Status = cluster.Status.Reset()
		cluster.Status.State = clusterregistryv1alpha1.ClusterStateDisabled
		currentConditions = ClusterConditionsMap{}
		err = UpdateCluster(ctx, reconcileError, r.GetClient(), cluster, currentConditions, log)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		if isClusterLocal {
			reconcileError = r.reconcileLocalCluster(ctx, cluster, currentConditions)
		} else {
			reconcileError = r.reconcileRemoteCluster(ctx, cluster, currentConditions)
		}

		SetCondition(cluster, currentConditions, ClusterReadyCondition(err), r.GetRecorder())
	}

	if isClusterLocal || reconcileError != nil {
		// status needs to be updated if the cluster is local or if there was any error setting up a remote cluster instance
		err = UpdateCluster(ctx, reconcileError, r.GetClient(), cluster, currentConditions, log)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info("reconciled successfully")

	return ctrl.Result{
		RequeueAfter: time.Second * time.Duration(r.config.ClusterController.RefreshIntervalSeconds),
	}, nil
}

func (r *ClusterReconciler) getK8SConfigForCluster(ctx context.Context, namespace string, name string) ([]byte, error) {
	var secret corev1.Secret
	err := r.GetClient().Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, &secret)
	if err != nil {
		return nil, err
	}

	if secret.Type != clusterregistryv1alpha1.SecretTypeClusterRegistry {
		return nil, errors.WithDetails(ErrInvalidSecret, "type", secret.Type)
	}

	if content, ok := secret.Data[clusterregistryv1alpha1.KubeconfigKey]; ok {
		return content, nil
	}

	return nil, ErrInvalidSecretContent
}

func (r *ClusterReconciler) getRemoteCluster(ctx context.Context, cluster *clusterregistryv1alpha1.Cluster) (*clusters.Cluster, error) {
	log := r.GetLogger().WithValues("cluster", cluster.Name)

	secretID := fmt.Sprintf("%s/%s", cluster.Spec.AuthInfo.SecretRef.Namespace, cluster.Spec.AuthInfo.SecretRef.Name)
	remoteCluster, _ := r.clustersManager.Get(cluster.Name)

	k8sconfig, err := r.getK8SConfigForCluster(ctx, cluster.Spec.AuthInfo.SecretRef.Namespace, cluster.Spec.AuthInfo.SecretRef.Name)
	if err != nil {
		cluster.Status.State = clusterregistryv1alpha1.ClusterStateInvalidAuthInfo

		if remoteCluster != nil {
			log.Info("cluster secret is removed")
			err := r.clustersManager.Remove(remoteCluster)
			if err != nil {
				return nil, errors.WrapIf(err, "could not remove cluster from manager")
			}
		}

		return nil, WrapAsPermanentError(err)
	}

	if remoteCluster != nil { // nolint:nestif
		if !remoteCluster.IsAlive() {
			return nil, WrapAsPermanentError(errors.New("remote cluster is not alive"))
		}
		credentialsChanged := false
		if remoteCluster.GetSecretID() != nil && *remoteCluster.GetSecretID() != secretID {
			log.Info("cluster secret reference changed")
			credentialsChanged = true
		}
		if len(remoteCluster.GetKubeconfig()) > 0 && !bytes.Equal(k8sconfig, remoteCluster.GetKubeconfig()) {
			log.Info("cluster secret content changed")
			credentialsChanged = true
		}
		if !credentialsChanged {
			return remoteCluster, nil
		}
		err := r.clustersManager.Remove(remoteCluster)
		if err != nil {
			return nil, errors.WrapIf(err, "could not remove cluster from manager")
		}
	}

	clusterConfig, err := clientcmd.Load(k8sconfig)
	if err != nil {
		return nil, errors.WrapIf(err, "could not load kubeconfig")
	}

	kubeConfigOverrides, err := util.GetKubeconfigOverridesForClusterByNetwork(cluster, r.config.NetworkName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rest, err := clientcmd.NewDefaultClientConfig(*clusterConfig, kubeConfigOverrides).ClientConfig()
	if err != nil {
		return nil, errors.WrapIf(err, "could not create k8s rest config")
	}

	onDeadFunc := func(c *clusters.Cluster) error {
		if r.queue != nil {
			r.queue.Add(reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: c.GetName(),
				},
			})
		}

		return nil
	}

	remoteCluster, err = clusters.NewCluster(ctx, cluster.Name, rest, r.GetLogger(), clusters.WithSecretID(secretID), clusters.WithCtrlOption(ctrl.Options{
		Scheme:             r.GetManager().GetScheme(),
		MetricsBindAddress: "0",
		Port:               0,
	}), clusters.WithOnDeadFunc(onDeadFunc), clusters.WithKubeconfig(k8sconfig))
	if err != nil {
		return nil, errors.WrapIf(err, "could not create new cluster")
	}

	err = remoteCluster.AddController(clusters.NewManagedController("remote-cluster", NewRemoteClusterReconciler(cluster.Name, r.GetManager(), r.GetLogger()), r.GetLogger()))
	if err != nil {
		return nil, errors.WrapIf(err, "could not add managed controller")
	}

	err = remoteCluster.AddController(clusters.NewManagedController("remote-cluster-feature", NewClusterFeatureReconciler(cluster.Name, remoteCluster, r.GetLogger()), r.GetLogger()))
	if err != nil {
		return nil, errors.WrapIf(err, "could not add managed controller")
	}

	err = remoteCluster.Start()
	if err != nil {
		return nil, errors.WrapIf(err, "could not start cluster")
	}

	err = r.clustersManager.Add(remoteCluster)
	if err != nil {
		return nil, errors.WrapIf(err, "could not add cluster to manager")
	}

	return remoteCluster, err
}

func (r *ClusterReconciler) removeRemoteCluster(name string) error {
	cluster, err := r.clustersManager.Get(name)
	if err == nil {
		err = r.clustersManager.Remove(cluster)
		if err != nil {
			return errors.WrapIf(err, "could not remove cluster from manager")
		}
	}

	return err
}

func (r *ClusterReconciler) reconcileRemoteCluster(ctx context.Context, cluster *clusterregistryv1alpha1.Cluster, currentConditions ClusterConditionsMap) error {
	SetCondition(cluster, currentConditions, LocalClusterCondition(false), r.GetRecorder())

	if cluster.Spec.AuthInfo.SecretRef.Name == "" || cluster.Spec.AuthInfo.SecretRef.Namespace == "" {
		cluster.Status.State = clusterregistryv1alpha1.ClusterStateMissingAuthInfo

		return nil
	}

	_, err := r.getRemoteCluster(ctx, cluster)
	if err != nil {
		err = errors.WithStackIf(err)
		SetCondition(cluster, currentConditions, ClustersSyncedCondition(err), r.GetRecorder())

		return errors.WithStackIf(err)
	}

	return nil
}

func (r *ClusterReconciler) reconcileLocalCluster(ctx context.Context, cluster *clusterregistryv1alpha1.Cluster, currentConditions ClusterConditionsMap) error {
	SetCondition(cluster, currentConditions, LocalClusterCondition(true), r.GetRecorder())

	err := r.checkLocalClusterConflict(ctx, cluster, currentConditions)
	if err != nil {
		return errors.WithStackIf(err)
	}

	clusterMetadata, err := clustermeta.GetClusterMetadata(ctx, r.GetClient())
	SetCondition(cluster, currentConditions, ClusterMetadataCondition(err), r.GetRecorder())
	if err != nil {
		return errors.WithStackIf(err)
	}
	cluster.Status.ClusterMetadata = clusterMetadata

	err = r.provisionLocalClusterReaderSecret(ctx, cluster)
	if err != nil {
		return errors.WithStackIf(err)
	}

	err = r.reconcileCoreSyncers()
	if err != nil {
		return errors.WithStackIf(err)
	}

	return nil
}

func (r *ClusterReconciler) provisionLocalClusterReaderSecret(ctx context.Context, cluster *clusterregistryv1alpha1.Cluster) error {
	if !r.config.ManageLocalClusterSecret {
		return nil
	}

	if cluster.Spec.AuthInfo.SecretRef.Name == "" || cluster.Spec.AuthInfo.SecretRef.Namespace == "" {
		return nil
	}

	secret, err := util.GetReaderSecretForCluster(ctx, r.GetClient(), r.GetManager().GetConfig(), cluster, types.NamespacedName{
		Name:      r.config.ReaderServiceAccountName,
		Namespace: r.config.Namespace,
	}, r.config.APIServerEndpointAddress)
	if err != nil {
		return errors.WithStackIf(err)
	}

	rec := reconciler.NewReconcilerWith(r.GetClient(),
		reconciler.WithLog(r.GetLogger()),
		reconciler.WithRecreateImmediately(),
		reconciler.WithEnableRecreateWorkload(),
		reconciler.WithRecreateEnabledForAll(),
	)
	_, err = rec.ReconcileResource(secret, reconciler.StatePresent)

	return errors.WithStackIf(err)
}

func (r *ClusterReconciler) checkLocalClusterConflict(ctx context.Context, cluster *clusterregistryv1alpha1.Cluster, currentConditions ClusterConditionsMap) error {
	clusters := &clusterregistryv1alpha1.ClusterList{}
	err := r.GetClient().List(ctx, clusters)
	if err != nil {
		return errors.WrapIf(err, "could not list objects")
	}

	for _, c := range clusters.Items {
		c := c
		if c.UID == cluster.UID {
			continue
		}
		if c.Spec.ClusterID == cluster.Spec.ClusterID {
			if GetCurrentCondition(&c, clusterregistryv1alpha1.ClusterConditionTypeLocalConflict).Status == corev1.ConditionFalse {
				r.triggerClusterReconcile(&c)
			}
			SetCondition(cluster, currentConditions, LocalClusterConflictCondition(true), r.GetRecorder())

			return WrapAsPermanentError(ErrLocalClusterConflict)
		}
	}

	SetCondition(cluster, currentConditions, LocalClusterConflictCondition(false), r.GetRecorder())

	return nil
}

func (r *ClusterReconciler) triggerClusterReconcile(cluster *clusterregistryv1alpha1.Cluster) {
	if r.queue != nil {
		r.queue.Add(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: cluster.GetName(),
			},
		})
	}
}

func (r *ClusterReconciler) triggerLocalClusterReconciles(ctx context.Context) error {
	clusters := &clusterregistryv1alpha1.ClusterList{}
	err := r.GetClient().List(ctx, clusters)
	if err != nil {
		return err
	}

	for _, c := range clusters.Items {
		c := c
		if GetCurrentCondition(&c, clusterregistryv1alpha1.ClusterConditionTypeLocalConflict).Status == corev1.ConditionTrue {
			r.triggerClusterReconcile(&c)
		}
	}

	return nil
}

func (r *ClusterReconciler) setClusterID(ctx context.Context) error {
	if r.clusterID == "" {
		clusterID, err := GetClusterID(ctx, r.GetClient())
		if err != nil {
			return errors.WrapIf(err, "could not get cluster id")
		}
		r.clusterID = clusterID
	}

	return nil
}

func (r *ClusterReconciler) setQueue(q workqueue.RateLimitingInterface) {
	r.queue = q
}

func (r *ClusterReconciler) SetupWithController(ctx context.Context, ctrl controller.Controller) error {
	err := r.ManagedReconciler.SetupWithController(ctx, ctrl)
	if err != nil {
		return errors.WithStack(err)
	}

	localClusterMapFunc := func(obj client.Object) []reconcile.Request {
		clusters := &clusterregistryv1alpha1.ClusterList{}
		err := r.GetClient().List(ctx, clusters)
		if err != nil {
			r.GetLogger().Error(err, "could not list clusters")

			return nil
		}

		items := []reconcile.Request{}

		for _, c := range clusters.Items {
			c := c
			if c.Status.Type == clusterregistryv1alpha1.ClusterTypeLocal {
				items = append(items, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      c.GetName(),
						Namespace: c.GetNamespace(),
					},
				})
			}
		}

		return items
	}

	for _, t := range []client.Object{
		&clusterregistryv1alpha1.ResourceSyncRule{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ResourceSyncRule",
				APIVersion: clusterregistryv1alpha1.SchemeBuilder.GroupVersion.String(),
			},
		},
		&clusterregistryv1alpha1.ResourceSyncRule{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterFeature",
				APIVersion: clusterregistryv1alpha1.SchemeBuilder.GroupVersion.String(),
			},
		},
	} {
		err = ctrl.Watch(&source.Kind{
			Type: t,
		}, handler.EnqueueRequestsFromMapFunc(localClusterMapFunc), &predicate.Funcs{
			UpdateFunc: func(e event.UpdateEvent) bool {
				if v, ok := e.ObjectOld.GetLabels()[CoreResourceLabelName]; ok && v == "true" {
					return true
				}

				return false
			},
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}

	err = ctrl.Watch(&InMemorySource{
		reconciler: r,
	}, handler.Funcs{})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *ClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	err := r.ManagedReconciler.SetupWithManager(ctx, mgr)
	if err != nil {
		return err
	}

	b := ctrl.NewControllerManagedBy(mgr)

	r.watchLocalClustersForConflict(ctx, b)
	r.watchClusterRegistrySecrets(ctx, b)

	ctrl, err := b.For(&clusterregistryv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: clusterregistryv1alpha1.SchemeBuilder.GroupVersion.String(),
		},
	}, builder.WithPredicates(predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectNew.GetGeneration() != e.ObjectOld.GetGeneration() {
				return true
			}
			if !reflect.DeepEqual(e.ObjectOld.GetAnnotations(), e.ObjectNew.GetAnnotations()) {
				return true
			}

			return false
		},
	})).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: r.config.ClusterController.WorkerCount,
		}).
		Build(r)
	if err != nil {
		return err
	}

	err = r.SetupWithController(ctx, ctrl)
	if err != nil {
		return err
	}

	r.SetClient(mgr.GetClient())

	return nil
}
