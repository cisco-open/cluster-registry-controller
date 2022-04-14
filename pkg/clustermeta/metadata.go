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

package clustermeta

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/cisco-open/cluster-registry-controller/api/v1alpha1"
)

func GetClusterMetadata(ctx context.Context, client client.Client) (v1alpha1.ClusterMetadata, error) {
	md := v1alpha1.ClusterMetadata{}

	nodes := &corev1.NodeList{}
	if err := client.List(ctx, nodes); err != nil {
		return md, err
	}

	if len(nodes.Items) == 0 {
		return md, nil
	}

	provider, err := DetectProvider(ctx, client, &nodes.Items[0])
	if err != nil && !IsUnknownProviderError(err) {
		return md, err
	}

	distribution, err := DetectDistribution(ctx, client, &nodes.Items[0])
	if err != nil && !IsUnknownDistributionError(err) {
		return md, err
	}

	md.Provider = provider
	md.Distribution = distribution

	kubeProxyVersions := make(map[string]struct{})
	kubeletVersions := make(map[string]struct{})
	regions := make(map[string]struct{})
	zones := make(map[string]struct{})

	for _, node := range nodes.Items {
		kubeProxyVersions[node.Status.NodeInfo.KubeProxyVersion] = struct{}{}
		kubeletVersions[node.Status.NodeInfo.KubeletVersion] = struct{}{}
		if len(node.Labels) > 0 {
			if v := node.Labels[corev1.LabelZoneRegionStable]; v != "" {
				regions[v] = struct{}{}
			}
			if v := node.Labels[corev1.LabelZoneFailureDomainStable]; v != "" {
				zones[v] = struct{}{}
			}
		}
	}

	if len(kubeProxyVersions) > 0 {
		for v := range kubeProxyVersions {
			md.KubeProxyVersions = append(md.KubeProxyVersions, v)
		}

		md.Version = md.KubeProxyVersions[0]
	}

	for v := range kubeletVersions {
		md.KubeletVersions = append(md.KubeletVersions, v)
	}

	if len(regions) > 0 || len(zones) > 0 {
		md.Locality = &v1alpha1.Locality{
			Regions: []string{},
			Zones:   []string{},
		}
	}

	if len(regions) > 0 {
		for v := range regions {
			md.Locality.Regions = append(md.Locality.Regions, v)
		}
		md.Locality.Region = strings.Join(md.Locality.Regions, ", ")
	}

	for v := range zones {
		md.Locality.Zones = append(md.Locality.Zones, v)
	}

	return md, nil
}
