// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package controllers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/resources"
)

const (
	CoreResourcesSourceFeatureName = "cluster-registry-core-resources-source"
	CoreResourceLabelName          = "cluster-registry-controller.k8s.cisco.com/core-sync-resource"
)

func (r *ClusterReconciler) reconcileCoreSyncers() error {
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

	for o, ds := range map[runtime.Object]reconciler.DesiredState{
		r.coreSyncersClusterFeature(): clusterFeatureDesiredState,
		r.clustersSyncRule():          reconciler.StatePresent,
		r.clusterSecretsSyncRule():    reconciler.StatePresent,
		r.resourceSyncRuleSync():      reconciler.StatePresent,
	} {
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
