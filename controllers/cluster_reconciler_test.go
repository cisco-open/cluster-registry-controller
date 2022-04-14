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

package controllers_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	clusterregistryv1alpha1 "github.com/cisco-open/cluster-registry-controller/api/v1alpha1"
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
