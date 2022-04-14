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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterregistryv1alpha1 "github.com/cisco-open/cluster-registry-controller/api/v1alpha1"
)

func (r *ClusterReconciler) watchLocalClustersForConflict(ctx context.Context, b *builder.Builder) {
	b.Watches(
		&source.Kind{Type: &clusterregistryv1alpha1.Cluster{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Cluster",
				APIVersion: clusterregistryv1alpha1.SchemeBuilder.GroupVersion.String(),
			},
		}},
		handler.EnqueueRequestsFromMapFunc(func(object client.Object) []ctrl.Request {
			reqs := make([]reconcile.Request, 0)
			clusters, err := GetClusters(ctx, r.GetClient())
			if err != nil {
				r.GetLogger().Error(err, "")

				return nil
			}

			for _, c := range clusters {
				nsn := types.NamespacedName{
					Name:      c.Name,
					Namespace: c.Namespace,
				}
				if c.Status.Type == clusterregistryv1alpha1.ClusterTypeLocal {
					reqs = append(reqs, ctrl.Request{
						NamespacedName: nsn,
					})
				}
			}

			return reqs
		}),
		builder.WithPredicates(&predicate.Funcs{
			GenericFunc: func(event.GenericEvent) bool { return false },
			UpdateFunc:  func(event.UpdateEvent) bool { return false },
		}))
}

func (r *ClusterReconciler) watchClusterRegistrySecrets(ctx context.Context, b *builder.Builder) {
	b.Watches(
		&source.Kind{Type: &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
		}},
		handler.EnqueueRequestsFromMapFunc(func(object client.Object) []ctrl.Request {
			reqs := make([]reconcile.Request, 0)
			if secret, ok := object.(*corev1.Secret); ok {
				if secret.Type != clusterregistryv1alpha1.SecretTypeClusterRegistry {
					return nil
				}
				clusters, err := GetClusters(ctx, r.GetClient())
				if err != nil {
					r.GetLogger().Error(err, "")

					return nil
				}

				for _, c := range clusters {
					nsn := types.NamespacedName{
						Name:      c.Name,
						Namespace: c.Namespace,
					}
					if c.Spec.AuthInfo.SecretRef.Name == secret.Name && c.Spec.AuthInfo.SecretRef.Namespace == secret.Namespace {
						reqs = append(reqs, ctrl.Request{
							NamespacedName: nsn,
						})
					}
				}

				return reqs
			}

			return nil
		}),
	)
}
