// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package controllers

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
)

type QueueAwareReconciler interface {
	setQueue(q workqueue.RateLimitingInterface)
}

type InMemorySource struct {
	reconciler QueueAwareReconciler
}

func (s *InMemorySource) String() string {
	return "in-memory"
}

func (s *InMemorySource) Start(ctx context.Context, h handler.EventHandler, q workqueue.RateLimitingInterface, p ...predicate.Predicate) error {
	s.reconciler.setQueue(q)

	return nil
}

func GetClusters(ctx context.Context, c client.Client) (map[types.UID]clusterregistryv1alpha1.Cluster, error) {
	clusterList := &clusterregistryv1alpha1.ClusterList{}
	err := c.List(ctx, clusterList)
	if err != nil {
		return nil, errors.WrapIf(err, "could not list clusters")
	}

	clusters := make(map[types.UID]clusterregistryv1alpha1.Cluster)
	for _, c := range clusterList.Items {
		clusters[c.Spec.ClusterID] = c
	}

	return clusters, nil
}

func GetClusterID(ctx context.Context, c client.Client) (types.UID, error) {
	ns := &corev1.Namespace{}
	err := c.Get(ctx, client.ObjectKey{
		Name: metav1.NamespaceSystem,
	}, ns)
	if err != nil {
		return "", err
	}

	return ns.UID, nil
}

func UpdateCluster(ctx context.Context, reconcileError error, c client.Client, cluster *clusterregistryv1alpha1.Cluster, currentConditions ClusterConditionsMap, log logr.Logger) error {
	conditions := make([]clusterregistryv1alpha1.ClusterCondition, 0)
	for _, condition := range currentConditions {
		// remove sync condition from local clusters
		if cluster.Status.Type == clusterregistryv1alpha1.ClusterTypeLocal && condition.Type == clusterregistryv1alpha1.ClusterConditionTypeClustersSynced {
			continue
		}
		conditions = append(conditions, condition)
	}
	cluster.Status.Conditions = conditions

	if reconcileError != nil {
		// status could be overwritten earlier, only override the ready state here
		if cluster.Status.State == clusterregistryv1alpha1.ClusterStateReady {
			cluster.Status.State = clusterregistryv1alpha1.ClusterStateFailed
		}
		cluster.Status.Message = reconcileError.Error()
		updateErr := UpdateClusterStatus(ctx, c, cluster, log)
		if updateErr != nil {
			log.Error(updateErr, "could not update resource status")
		}
		//nolint:errorlint
		if e, ok := reconcileError.(interface{ IsPermanent() bool }); ok && e.IsPermanent() {
			log.Error(reconcileError, "", errors.GetDetails(reconcileError)...)
			reconcileError = nil
		}

		return reconcileError
	}

	err := UpdateClusterStatus(ctx, c, cluster, log)
	if err != nil {
		return errors.WithStackIf(err)
	}

	return nil
}

func UpdateClusterStatus(ctx context.Context, c client.Client, cluster *clusterregistryv1alpha1.Cluster, log logr.Logger) error {
	desired := cluster.DeepCopy()

	log.Info("update cluster status")

	err := c.Status().Update(ctx, cluster)
	if apierrors.IsConflict(err) {
		current := cluster.DeepCopy()
		err = c.Get(ctx, client.ObjectKey{
			Name: current.Name,
		}, current)
		if err != nil {
			return errors.WrapIf(err, "could not get cluster")
		}
		desired.SetResourceVersion(current.GetResourceVersion())
		err = c.Status().Update(ctx, desired)
		if err != nil {
			return errors.WrapIf(err, "could not update cluster status")
		}
	}
	if err != nil {
		return errors.WrapIf(err, "could not update status")
	}

	return nil
}
