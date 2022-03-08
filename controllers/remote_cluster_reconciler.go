// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package controllers

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/clustermeta"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/clusters"
)

type RemoteClusterReconciler struct {
	clusters.ManagedReconciler

	localMgr      ctrl.Manager
	localRecorder record.EventRecorder

	clusterID types.UID
}

func NewRemoteClusterReconciler(name string, localMgr ctrl.Manager, log logr.Logger) *RemoteClusterReconciler {
	return &RemoteClusterReconciler{
		ManagedReconciler: clusters.NewManagedReconciler(name, log),

		localMgr:      localMgr,
		localRecorder: localMgr.GetEventRecorderFor("cluster-controller"),
	}
}

func (r *RemoteClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.GetLogger().WithValues("cluster", req.NamespacedName)

	err := r.setClusterID(ctx)
	if err != nil {
		return ctrl.Result{}, errors.WithStackIf(err)
	}

	log.Info("reconciling")

	cluster := &clusterregistryv1alpha1.Cluster{}
	err = r.localMgr.GetClient().Get(ctx, req.NamespacedName, cluster)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "could not get object")
	}

	cluster.Status = cluster.Status.Reset()
	currentConditions := GetCurrentConditions(cluster)

	reconcileError := r.reconcile(ctx, cluster, currentConditions)

	if reconcileError != nil {
		SetCondition(cluster, currentConditions, ClustersSyncedCondition(reconcileError), r.localRecorder)
	}

	log.Info("update cluster status")
	err = UpdateCluster(ctx, reconcileError, r.localMgr.GetClient(), cluster, currentConditions, r.GetLogger())
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled successfully", "uid", cluster.UID)

	return ctrl.Result{}, nil
}

func (r *RemoteClusterReconciler) reconcile(ctx context.Context, cluster *clusterregistryv1alpha1.Cluster, currentConditions ClusterConditionsMap) error {
	if r.clusterID != cluster.Spec.ClusterID {
		return WrapAsPermanentError(errors.WithDetails(ErrInvalidClusterID, "nsUID", r.clusterID, "clusterUID", cluster.UID))
	}

	clusterMetadata, err := clustermeta.GetClusterMetadata(ctx, r.GetClient())
	SetCondition(cluster, currentConditions, ClusterMetadataCondition(err), r.localRecorder)
	if err != nil {
		err = errors.WithStackIf(err)
		SetCondition(cluster, currentConditions, ClustersSyncedCondition(err), r.localRecorder)

		return errors.WithStackIf(err)
	}

	cluster.Status.ClusterMetadata = clusterMetadata

	SetCondition(cluster, currentConditions, ClustersSyncedCondition(r.syncClusters(ctx)), r.localRecorder)

	return nil
}

func (r *RemoteClusterReconciler) setClusterID(ctx context.Context) error {
	if r.clusterID == "" {
		clusterID, err := GetClusterID(ctx, r.GetClient())
		if err != nil {
			return errors.WrapIf(err, "could not get cluster id")
		}
		r.clusterID = clusterID
	}

	return nil
}

func (r *RemoteClusterReconciler) syncClusters(ctx context.Context) error {
	localClusters, err := GetClusters(ctx, r.localMgr.GetClient())
	if err != nil {
		return errors.WrapIf(err, "could not get local clusters")
	}

	remoteClusters, err := GetClusters(ctx, r.GetClient())
	if err != nil {
		return errors.WrapIf(err, "could not get remote clusters")
	}

	var errs error
	for _, c := range localClusters {
		if rc, ok := remoteClusters[c.Spec.ClusterID]; ok {
			if c.Name != rc.Name {
				errs = errors.Append(errs, errors.WithDetails(errors.Errorf("cluster name mismatch %s != %s", c.Name, rc.Name), "localName", c.Name, "remoteName", rc.Name))
			}

			if rc.Status.State != clusterregistryv1alpha1.ClusterStateReady {
				errs = errors.Append(errs, errors.WithDetails(errors.Errorf("cluster '%s' status is not ready", rc.Name), "status", rc.Status.State))
			}
		} else {
			errs = errors.Append(errs, errors.WithDetails(errors.Errorf("cluster '%s' not found at remote side", c.Name), "uid", c.Spec.ClusterID, "name", c.Name))
		}
	}

	return errs
}

func (r *RemoteClusterReconciler) SetupWithController(ctx context.Context, ctrl controller.Controller) error {
	err := r.ManagedReconciler.SetupWithController(ctx, ctrl)
	if err != nil {
		return err
	}

	err = ctrl.Watch(
		&source.Kind{
			Type: &clusterregistryv1alpha1.Cluster{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Cluster",
					APIVersion: clusterregistryv1alpha1.SchemeBuilder.GroupVersion.String(),
				},
			},
		},
		handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
			return []reconcile.Request{
				{
					NamespacedName: client.ObjectKey{
						Name: r.GetName(),
					},
				},
			}
		}),
		predicate.Funcs{ // reconcile the local peer cluster for any change at the remote side
			UpdateFunc: func(e event.UpdateEvent) bool {
				var ok bool
				var oldObject, newObject *clusterregistryv1alpha1.Cluster

				if oldObject, ok = e.ObjectOld.(*clusterregistryv1alpha1.Cluster); !ok {
					return false
				}
				if newObject, ok = e.ObjectNew.(*clusterregistryv1alpha1.Cluster); !ok {
					return false
				}

				if oldObject.Status.State != newObject.Status.State || oldObject.Status.Message != newObject.Status.Message {
					return true
				}

				return false
			},
		})
	if err != nil {
		return errors.WrapIf(err, "could not create watch for remote clusters")
	}

	clusterInformer, err := r.localMgr.GetCache().GetInformerForKind(ctx, clusterregistryv1alpha1.SchemeBuilder.GroupVersion.WithKind("Cluster"))
	if err != nil {
		return errors.WrapIf(err, "could not create local informer for clusters")
	}

	err = ctrl.Watch(
		&source.Informer{
			Informer: clusterInformer,
		},
		handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
			return []reconcile.Request{
				{
					NamespacedName: client.ObjectKey{
						Name: r.GetName(),
					},
				},
			}
		}),
		predicate.GenerationChangedPredicate{})
	if err != nil {
		return errors.WrapIf(err, "could not create watch for local informer")
	}

	return nil
}
