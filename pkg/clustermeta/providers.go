// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

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
)

func IsAmazon(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	ok, err := detectNodeByProviderID(ctx, client, node, "aws")
	if err != nil {
		return ok, "", err
	}

	return ok, AMAZON, nil
}

func IsAzure(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	ok, err := detectNodeByProviderID(ctx, client, node, "azure")
	if err != nil {
		return ok, "", err
	}

	return ok, AZURE, nil
}

func IsGoogle(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	ok, err := detectNodeByProviderID(ctx, client, node, "gce")
	if err != nil {
		return ok, "", err
	}

	return ok, GOOGLE, nil
}

func IsVsphere(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	ok, err := detectNodeByProviderID(ctx, client, node, "vsphere")
	if err != nil {
		return ok, "", err
	}

	return ok, VSPHERE, nil
}

func IsKind(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error) {
	ok, err := detectNodeByProviderID(ctx, client, node, "kind")
	if err != nil {
		return ok, "", err
	}

	return ok, KINDP, nil
}
