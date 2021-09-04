// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

package controllers

import (
	"context"
	"fmt"
	"net/url"
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
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/internal/config"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/clustermeta"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/clusters"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/util"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
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

func (r *ClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	log := r.GetLogger().WithValues("cluster", req.NamespacedName)

	err := r.setClusterID()
	if err != nil {
		return ctrl.Result{}, errors.WithStackIf(err)
	}

	log.Info("reconciling")

	cluster := &clusterregistryv1alpha1.Cluster{}
	err = r.GetManager().GetClient().Get(r.GetContext(), req.NamespacedName, cluster)
	if apierrors.IsNotFound(err) {
		removeErr := r.removeRemoteCluster(req.NamespacedName.Name)
		if removeErr != nil && !errors.Is(errors.Cause(removeErr), clusters.ErrClusterNotFound) {
			return ctrl.Result{}, errors.WithStackIf(removeErr)
		}

		// trigger local reconciles to fix possible local conflict status
		if err = r.triggerLocalClusterReconciles(); err != nil {
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
		if err := r.triggerLocalClusterReconciles(); err != nil {
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
		err = UpdateCluster(r.GetContext(), reconcileError, r.GetManager().GetClient(), cluster, currentConditions, log)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		if isClusterLocal {
			reconcileError = r.reconcileLocalCluster(r.GetContext(), cluster, currentConditions)
		} else {
			reconcileError = r.reconcileRemoteCluster(r.GetContext(), cluster, currentConditions)
		}

		SetCondition(cluster, currentConditions, ClusterReadyCondition(err), r.GetRecorder())
	}

	if isClusterLocal || reconcileError != nil {
		// status needs to be updated if the cluster is local or if there was any error setting up a remote cluster instance
		err = UpdateCluster(r.GetContext(), reconcileError, r.GetManager().GetClient(), cluster, currentConditions, log)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info("reconciled successfully")

	return ctrl.Result{
		RequeueAfter: time.Second * time.Duration(r.config.ClusterController.RefreshIntervalSeconds),
	}, nil
}

func (r *ClusterReconciler) getK8SConfigForCluster(namespace string, name string) ([]byte, error) {
	var secret corev1.Secret
	err := r.GetManager().GetClient().Get(context.TODO(), client.ObjectKey{
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
	secretID := fmt.Sprintf("%s/%s", cluster.Spec.AuthInfo.SecretRef.Namespace, cluster.Spec.AuthInfo.SecretRef.Name)

	remoteCluster, _ := r.clustersManager.Get(cluster.Name)
	if remoteCluster != nil {
		if !remoteCluster.IsAlive() {
			return nil, WrapAsPermanentError(errors.New("remote cluster is not alive"))
		}
		if remoteCluster.GetSecretID() == nil || *remoteCluster.GetSecretID() == secretID {
			return remoteCluster, nil
		}
		err := r.clustersManager.Remove(remoteCluster)
		if err != nil {
			return nil, errors.WrapIf(err, "could not remove cluster from manager")
		}
	}

	k8sconfig, err := r.getK8SConfigForCluster(cluster.Spec.AuthInfo.SecretRef.Namespace, cluster.Spec.AuthInfo.SecretRef.Name)
	if err != nil {
		cluster.Status.State = clusterregistryv1alpha1.ClusterStateInvalidAuthInfo

		return nil, WrapAsPermanentError(err)
	}

	clusterConfig, err := clientcmd.Load(k8sconfig)
	if err != nil {
		return nil, errors.WrapIf(err, "could not load kubeconfig")
	}

	rest, err := clientcmd.NewDefaultClientConfig(*clusterConfig, r.getKubeconfigOverridesForNetwork(cluster, r.config.NetworkName)).ClientConfig()
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
	}), clusters.WithOnDeadFunc(onDeadFunc))
	if err != nil {
		return nil, errors.WrapIf(err, "could not create new cluster")
	}

	err = remoteCluster.AddController(clusters.NewManagedController("remote-cluster", NewRemoteClusterReconciler(cluster.Name, r.GetManager(), r.GetLogger()), r.GetLogger()))
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

func (r *ClusterReconciler) getKubeconfigOverridesForNetwork(cluster *clusterregistryv1alpha1.Cluster, networkName string) *clientcmd.ConfigOverrides {
	overrides := &clientcmd.ConfigOverrides{}

	if len(cluster.Spec.KubernetesAPIEndpoints) == 0 {
		return overrides
	}

	var endpoint clusterregistryv1alpha1.KubernetesAPIEndpoint
	for _, apiEndpoint := range cluster.Spec.KubernetesAPIEndpoints {
		if apiEndpoint.ClientNetwork == networkName {
			endpoint = apiEndpoint
			break
		}
		// use for every network if the endpoint is not network specific
		if apiEndpoint.ClientNetwork == "" {
			endpoint = apiEndpoint
		}
	}

	if endpoint.ServerAddress != "" {
		overrides.ClusterInfo.Server = (&url.URL{
			Scheme: "https",
			Host:   endpoint.ServerAddress,
		}).String()
	}

	if len(endpoint.CABundle) > 0 {
		overrides.ClusterInfo.CertificateAuthorityData = endpoint.CABundle
	}

	return overrides
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

	clusterMetadata, err := clustermeta.GetClusterMetadata(ctx, r.GetManager().GetClient())
	SetCondition(cluster, currentConditions, ClusterMetadataCondition(err), r.GetRecorder())
	if err != nil {
		return errors.WithStackIf(err)
	}
	cluster.Status.ClusterMetadata = clusterMetadata

	err = r.provisionLocalClusterReaderSecret(ctx, cluster)
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

	secret, err := util.GetReaderSecretForCluster(ctx, r.GetManager().GetClient(), r.GetManager().GetConfig(), cluster, types.NamespacedName{
		Name:      r.config.ReaderServiceAccountName,
		Namespace: r.config.Namespace,
	})
	if err != nil {
		return errors.WithStackIf(err)
	}

	rec := reconciler.NewReconcilerWith(r.GetManager().GetClient(),
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
	err := r.GetManager().GetClient().List(ctx, clusters)
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

func (r *ClusterReconciler) triggerLocalClusterReconciles() error {
	clusters := &clusterregistryv1alpha1.ClusterList{}
	err := r.GetManager().GetClient().List(r.GetContext(), clusters)
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

func (r *ClusterReconciler) setClusterID() error {
	if r.clusterID == "" {
		clusterID, err := GetClusterID(r.GetContext(), r.GetManager().GetClient())
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
		return err
	}

	return ctrl.Watch(&InMemorySource{
		reconciler: r,
	}, handler.Funcs{})
}

func (r *ClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	err := r.ManagedReconciler.SetupWithManager(ctx, mgr)
	if err != nil {
		return err
	}

	b := ctrl.NewControllerManagedBy(mgr)

	r.watchLocalClustersForConflict(b)
	r.watchClusterRegistrySecrets(b)

	ctrl, err := b.For(&clusterregistryv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: clusterregistryv1alpha1.SchemeBuilder.GroupVersion.String(),
		},
	}, builder.WithPredicates(predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.MetaNew.GetGeneration() != e.MetaOld.GetGeneration() {
				return true
			}
			if !reflect.DeepEqual(e.MetaOld.GetAnnotations(), e.MetaNew.GetAnnotations()) {
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

	return nil
}
