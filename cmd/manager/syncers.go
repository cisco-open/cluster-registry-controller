// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

package main

import (
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/controllers"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/internal/config"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/clusters"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
	"github.com/banzaicloud/operator-tools/pkg/resources"
)

func AddClustersSyncRule(clustersManager *clusters.Manager, mgr ctrl.Manager, log logr.Logger, config config.Configuration) {
	clustersManager.AddOnAfterAddFunc(func(c *clusters.Cluster) {
		sr := &clusterregistryv1alpha1.ResourceSyncRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: "sync-clusters",
			},
			Spec: clusterregistryv1alpha1.ResourceSyncRuleSpec{
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

		_, err := controllers.InitNewResourceSyncController(sr, c, clustersManager, mgr, ctrl.Log, config)
		if err != nil {
			log.Error(err, "could not initialize sync controller")
		}
	})
}

func AddClusterSecretsSyncRule(clustersManager *clusters.Manager, mgr ctrl.Manager, log logr.Logger, config config.Configuration) {
	clustersManager.AddOnAfterAddFunc(func(c *clusters.Cluster) {
		sr := &clusterregistryv1alpha1.ResourceSyncRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: "sync-cluster-secrets",
			},
			Spec: clusterregistryv1alpha1.ResourceSyncRuleSpec{
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
									config.Namespace,
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

		_, err := controllers.InitNewResourceSyncController(sr, c, clustersManager, mgr, ctrl.Log, config)
		if err != nil {
			log.Error(err, "could not initialize sync controller")
		}
	})
}

func AddResourceSyncRuleSyncRule(clustersManager *clusters.Manager, mgr ctrl.Manager, log logr.Logger, config config.Configuration) {
	clustersManager.AddOnAfterAddFunc(func(c *clusters.Cluster) {
		sr := &clusterregistryv1alpha1.ResourceSyncRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: "sync-resource-sync-rules",
			},
			Spec: clusterregistryv1alpha1.ResourceSyncRuleSpec{
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

		_, err := controllers.InitNewResourceSyncController(sr, c, clustersManager, mgr, ctrl.Log, config)
		if err != nil {
			log.Error(err, "could not initialize sync controller")
		}
	})
}
