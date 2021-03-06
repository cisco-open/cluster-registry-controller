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

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	AMAZON  = "amazon"
	AZURE   = "azure"
	GOOGLE  = "google"
	VSPHERE = "vsphere"
	KINDP   = "kind"
	CISCO   = "cisco"
)

func IsAmazon(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	ok, _, err := detectNodeByProviderID(ctx, client, node, "aws")
	if err != nil {
		return ok, "", err
	}

	return ok, AMAZON, nil
}

func IsAzure(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	ok, _, err := detectNodeByProviderID(ctx, client, node, "azure")
	if err != nil {
		return ok, "", err
	}

	return ok, AZURE, nil
}

func IsGoogle(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	ok, _, err := detectNodeByProviderID(ctx, client, node, "gce")
	if err != nil {
		return ok, "", err
	}

	return ok, GOOGLE, nil
}

func IsVsphere(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	var ok bool
	var err error

	ok, node, err = detectNodeByProviderID(ctx, client, node, "vsphere")
	if err != nil {
		return ok, "", err
	}

	if node.Labels == nil {
		node.Labels = make(map[string]string)
	}

	if value, ok := node.Labels["iks.intersight.cisco.com/version"]; ok && value != "" {
		return false, "", nil
	}

	return ok, VSPHERE, nil
}

func IsKind(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	ok, _, err := detectNodeByProviderID(ctx, client, node, "kind")
	if err != nil {
		return ok, "", err
	}

	return ok, KINDP, nil
}

func IsCisco(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	ok, node, err := detectNodeByProviderID(ctx, client, node, "vsphere")
	if err != nil {
		return ok, "", err
	}

	if value, ok := node.Labels["iks.intersight.cisco.com/version"]; !ok || value == "" {
		return false, "", nil
	}

	return ok, CISCO, nil
}
