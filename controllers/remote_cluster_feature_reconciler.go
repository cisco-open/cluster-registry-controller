// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterregistryv1alpha1 "github.com/cisco-open/cluster-registry-controller/api/v1alpha1"
	"github.com/cisco-open/cluster-registry-controller/pkg/clusters"
)

type ClusterFeatureReconciler struct {
	clusters.ManagedReconciler

	cluster *clusters.Cluster
}

func NewClusterFeatureReconciler(name string, cluster *clusters.Cluster, log logr.Logger) *ClusterFeatureReconciler {
	return &ClusterFeatureReconciler{
		ManagedReconciler: clusters.NewManagedReconciler(name, log),

		cluster: cluster,
	}
}

func (r *ClusterFeatureReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.GetLogger().WithValues("clusterFeature", req.NamespacedName, "cluster", r.cluster.GetName())

	feature := &clusterregistryv1alpha1.ClusterFeature{}
	err := r.GetClient().Get(ctx, req.NamespacedName, feature)
	if apierrors.IsNotFound(err) {
		log.Info("delete feature")
		r.cluster.RemoveFeature(req.NamespacedName.String())

		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "could not get object")
	}

	log.Info("add feature")
	r.cluster.AddFeature(clusters.NewClusterFeature(req.NamespacedName.String(), feature.Spec.FeatureName, feature.GetLabels()))

	return ctrl.Result{}, nil
}

func (r *ClusterFeatureReconciler) SetupWithController(ctx context.Context, ctrl controller.Controller) error {
	err := r.ManagedReconciler.SetupWithController(ctx, ctrl)
	if err != nil {
		return err
	}

	err = ctrl.Watch(
		&source.Kind{
			Type: &clusterregistryv1alpha1.ClusterFeature{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterFeature",
					APIVersion: clusterregistryv1alpha1.SchemeBuilder.GroupVersion.String(),
				},
			},
		},
		&handler.EnqueueRequestForObject{},
		predicate.Funcs{},
	)
	if err != nil {
		return errors.WrapIf(err, "could not create watch for cluster features")
	}

	return nil
}
