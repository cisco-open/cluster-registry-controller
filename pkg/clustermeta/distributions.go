// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package clustermeta

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	PKE   = "PKE"
	EKS   = "EKS"
	GKE   = "GKE"
	AKS   = "AKS"
	KINDD = "KIND"
	IKS   = "IKS"
)

func IsPKE(ctx context.Context, client client.Client, node *corev1.Node) (match bool, distribution string, err error) {
	distribution = PKE

	if node == nil {
		node, _, err = getK8sNode(ctx, client)
		if err != nil {
			return
		}
	}

	match = true

	if value, ok := node.Labels["nodepool.banzaicloud.io/name"]; !ok || value == "" {
		match = false

		return
	}

	if value, ok := node.Annotations["kubeadm.alpha.kubernetes.io/cri-socket"]; !ok || value == "" {
		match = false

		return
	}

	return match, distribution, err
}

func IsEKS(ctx context.Context, client client.Client, node *corev1.Node) (match bool, distribution string, err error) {
	distribution = EKS

	if node == nil {
		node, _, err = getK8sNode(ctx, client)
		if err != nil {
			return
		}
	}

	match = true

	if _, ok := node.Annotations["kubeadm.alpha.kubernetes.io/cri-socket"]; ok {
		match = false

		return
	}

	var provider string

	provider, err = DetectProvider(ctx, client, node)
	if IsUnknownProviderError(err) {
		return false, distribution, nil
	}
	if err != nil {
		return
	}

	if provider != AMAZON {
		match = false

		return
	}

	return match, distribution, err
}

func IsGKE(ctx context.Context, client client.Client, node *corev1.Node) (match bool, distribution string, err error) {
	distribution = GKE

	if node == nil {
		node, _, err = getK8sNode(ctx, client)
		if err != nil {
			return
		}
	}

	match = true

	if value, ok := node.Labels["cloud.google.com/gke-nodepool"]; !ok || value == "" {
		match = false

		return
	}

	var provider string

	provider, err = DetectProvider(ctx, client, node)
	if IsUnknownProviderError(err) {
		return false, distribution, nil
	}
	if err != nil {
		return
	}

	if provider != GOOGLE {
		match = false

		return
	}

	return match, distribution, err
}

func IsAKS(ctx context.Context, client client.Client, node *corev1.Node) (match bool, distribution string, err error) {
	distribution = AKS

	if node == nil {
		node, _, err = getK8sNode(ctx, client)
		if err != nil {
			return
		}
	}

	match = true
	if value, ok := node.Labels["agentpool"]; !ok || value == "" {
		match = false

		return
	}

	var provider string
	provider, err = DetectProvider(ctx, client, node)
	if IsUnknownProviderError(err) {
		return false, distribution, nil
	}
	if err != nil {
		return
	}

	if provider != AZURE {
		match = false

		return
	}

	return match, distribution, err
}

func IsKIND(ctx context.Context, client client.Client, node *corev1.Node) (match bool, distribution string, err error) {
	distribution = KINDD

	if node == nil {
		node, _, err = getK8sNode(ctx, client)
		if err != nil {
			return
		}
	}

	match = true

	var provider string
	provider, err = DetectProvider(ctx, client, node)
	if IsUnknownProviderError(err) {
		return false, distribution, nil
	}
	if err != nil {
		return
	}

	if provider != KINDP {
		match = false

		return
	}

	return match, distribution, err
}

func IsIKS(ctx context.Context, client client.Client, node *corev1.Node) (match bool, distribution string, err error) {
	distribution = IKS

	if node == nil {
		node, _, err = getK8sNode(ctx, client)
		if err != nil {
			return
		}
	}

	match = true
	if value, ok := node.Labels["iks.intersight.cisco.com/version"]; !ok || value == "" {
		match = false

		return
	}

	var provider string
	provider, err = DetectProvider(ctx, client, node)
	if IsUnknownProviderError(err) {
		return false, distribution, nil
	}
	if err != nil {
		return
	}

	if provider != CISCO {
		match = false

		return
	}

	return match, distribution, err
}
