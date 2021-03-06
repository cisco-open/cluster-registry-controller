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
	"emperror.dev/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/resources"
	clusterregistryv1alpha1 "github.com/cisco-open/cluster-registry-controller/api/v1alpha1"
)

const (
	CoreResourcesSourceFeatureName = "cluster-registry-core-resources-source"
	CoreResourceLabelName          = "cluster-registry-controller.k8s.cisco.com/core-sync-resource"
)

func (r *ClusterReconciler) reconcileCoreSyncers(cluster *clusterregistryv1alpha1.Cluster) error {
	rec := reconciler.NewGenericReconciler(
		r.GetManager().GetClient(),
		r.GetLogger(),
		reconciler.ReconcilerOpts{
			EnableRecreateWorkloadOnImmutableFieldChange: true,
			Scheme: r.GetManager().GetScheme(),
		},
	)

	clusterFeatureDesiredState := reconciler.StateAbsent
	if r.config.CoreResourcesSourceEnabled {
		clusterFeatureDesiredState = reconciler.StatePresent
	}

	for o, ds := range map[client.Object]reconciler.DesiredState{
		r.coreSyncersClusterFeature(): clusterFeatureDesiredState,
		r.clustersSyncRule():          reconciler.StatePresent,
		r.clusterSecretsSyncRule():    reconciler.StatePresent,
		r.resourceSyncRuleSync():      reconciler.StatePresent,
	} {
		// set resource owner
		if err := controllerutil.SetControllerReference(cluster, o, r.GetManager().GetScheme()); err != nil {
			return errors.WrapIfWithDetails(err, "couldn't set local Cluster as resource owner for resource",
				"namespace", o.GetNamespace(), "name", o.GetName())
		}
		_, err := rec.ReconcileResource(o, ds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ClusterReconciler) coreSyncersClusterFeature() *clusterregistryv1alpha1.ClusterFeature {
	return &clusterregistryv1alpha1.ClusterFeature{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-registry-core-resources",
			Labels: map[string]string{
				CoreResourceLabelName: "true",
			},
		},
		Spec: clusterregistryv1alpha1.ClusterFeatureSpec{
			FeatureName: CoreResourcesSourceFeatureName,
		},
	}
}

func (r *ClusterReconciler) clustersSyncRule() *clusterregistryv1alpha1.ResourceSyncRule {
	return &clusterregistryv1alpha1.ResourceSyncRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-registry-core-resource-clusters-sink",
			Labels: map[string]string{
				CoreResourceLabelName: "true",
			},
		},
		Spec: clusterregistryv1alpha1.ResourceSyncRuleSpec{
			ClusterFeatureMatches: []clusterregistryv1alpha1.ClusterFeatureMatch{
				{
					FeatureName: CoreResourcesSourceFeatureName,
				},
			},
			GVK: resources.GroupVersionKind(clusterregistryv1alpha1.SchemeBuilder.GroupVersion.WithKind("Cluster")),
			Rules: []clusterregistryv1alpha1.SyncRule{
				{
					Mutations: clusterregistryv1alpha1.Mutations{
						SyncStatus: false,
					},
				},
			},
		},
	}
}

func (r *ClusterReconciler) clusterSecretsSyncRule() *clusterregistryv1alpha1.ResourceSyncRule {
	return &clusterregistryv1alpha1.ResourceSyncRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-registry-core-resource-cluster-secrets-sink",
			Labels: map[string]string{
				CoreResourceLabelName: "true",
			},
		},
		Spec: clusterregistryv1alpha1.ResourceSyncRuleSpec{
			ClusterFeatureMatches: []clusterregistryv1alpha1.ClusterFeatureMatch{
				{
					FeatureName: CoreResourcesSourceFeatureName,
				},
			},
			GVK: resources.GroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret")),
			Rules: []clusterregistryv1alpha1.SyncRule{
				{
					Matches: []clusterregistryv1alpha1.SyncRuleMatch{
						{
							Content: []clusterregistryv1alpha1.ContentSelector{
								{
									Key:   "type",
									Value: intstr.FromString(string(clusterregistryv1alpha1.SecretTypeClusterRegistry)),
								},
							},
							Namespaces: []string{
								r.config.Namespace,
							},
						},
					},
					Mutations: clusterregistryv1alpha1.Mutations{
						SyncStatus: false,
					},
				},
			},
		},
	}
}

func (r *ClusterReconciler) resourceSyncRuleSync() *clusterregistryv1alpha1.ResourceSyncRule {
	return &clusterregistryv1alpha1.ResourceSyncRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-registry-core-resource-syncrules-sink",
			Labels: map[string]string{
				CoreResourceLabelName: "true",
			},
		},
		Spec: clusterregistryv1alpha1.ResourceSyncRuleSpec{
			ClusterFeatureMatches: []clusterregistryv1alpha1.ClusterFeatureMatch{
				{
					FeatureName: CoreResourcesSourceFeatureName,
				},
			},
			GVK: resources.GroupVersionKind(clusterregistryv1alpha1.SchemeBuilder.GroupVersion.WithKind("ResourceSyncRule")),
			Rules: []clusterregistryv1alpha1.SyncRule{
				{
					Mutations: clusterregistryv1alpha1.Mutations{
						SyncStatus: false,
					},
				},
			},
		},
	}
}
