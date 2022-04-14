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

package util

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/cisco-open/cluster-registry-controller/api/v1alpha1"
	"github.com/cisco-open/cluster-registry-controller/internal/config"
)

func ProvisionLocalClusterObject(c client.Client, log logr.Logger, configuration config.Configuration) error {
	clusterSpec := v1alpha1.Cluster{}

	err := c.Get(context.Background(), types.NamespacedName{
		Name:      configuration.ProvisionLocalCluster,
		Namespace: configuration.Namespace,
	}, &clusterSpec)

	if err == nil {
		log.Info("local cluster object already exists, skipping provisioning", "cluster_name", configuration.ProvisionLocalCluster)

		return nil
	}

	newClusterSpec, err := NewLocalCluster(c, configuration.Namespace, configuration.ProvisionLocalCluster, configuration.APIServerEndpointAddress)
	if err != nil {
		return errors.WrapIf(err, "cannot create new local cluster object")
	}

	err = c.Create(context.Background(), newClusterSpec)
	if err != nil {
		return errors.WrapIf(err, "cannot create new local cluster object")
	}

	log.Info("provisioned local cluster configuration", "cluster_name", configuration.ProvisionLocalCluster)

	return nil
}

func NewLocalCluster(c client.Client, namespace, name, apiServerEndpointAddress string) (*v1alpha1.Cluster, error) {
	ns := &corev1.Namespace{}
	err := c.Get(context.Background(), types.NamespacedName{
		Name: metav1.NamespaceSystem,
	}, ns)
	if err != nil {
		return nil, err
	}

	return &v1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1alpha1.ClusterSpec{
			ClusterID: ns.UID,
			AuthInfo: v1alpha1.AuthInfo{
				SecretRef: v1alpha1.NamespacedName{
					Name:      name,
					Namespace: namespace,
				},
			},
			KubernetesAPIEndpoints: []v1alpha1.KubernetesAPIEndpoint{
				{
					ServerAddress: apiServerEndpointAddress,
				},
			},
		},
	}, nil
}
