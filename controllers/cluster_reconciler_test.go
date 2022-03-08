// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package controllers_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
)

var _ = Describe("Cluster controller", func() {
	const (
		ClusterName = "demo"

		timeout = time.Second * 10
		// duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	It("reconciles objects properly", func() {
		By("By creating a new Cluster")
		ctx := context.Background()
		cluster := &clusterregistryv1alpha1.Cluster{
			TypeMeta: metav1.TypeMeta{
				APIVersion: clusterregistryv1alpha1.GroupVersion.String(),
				Kind:       "Cluster",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: ClusterName,
			},
			Spec: clusterregistryv1alpha1.ClusterSpec{
				ClusterID: "abc",
			},
		}
		Expect(k8sClient.Create(ctx, cluster)).Should(Succeed())

		createdCluster := &clusterregistryv1alpha1.Cluster{}
		Eventually(func() bool {
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name: cluster.Name,
			}, createdCluster)

			return err == nil
		}, timeout, interval).Should(BeTrue())
		Expect(createdCluster.Spec.ClusterID).Should(Equal(cluster.Spec.ClusterID))
	})
})
