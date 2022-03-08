// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package clustermeta

import (
	"context"
	"errors"
	"net/url"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var knownProviders = []IsProvider{
	IsAmazon,
	IsAzure,
	IsGoogle,
	IsVsphere,
	IsKind,
	IsCisco,
}

type IsProvider func(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error)

type UnknownProviderError struct{}

func (UnknownProviderError) Error() string {
	return "unknown provider"
}

func IsUnknownProviderError(err error) bool {
	return errors.As(err, &UnknownProviderError{})
}

func DetectProvider(ctx context.Context, client client.Client, node *corev1.Node) (string, error) {
	for _, f := range knownProviders {
		select {
		case <-ctx.Done():
			return "", UnknownProviderError{}
		default:
			if ok, providerName, err := f(ctx, client, node); err != nil {
				return "", err
			} else if ok {
				return providerName, nil
			}
		}
	}

	return "", UnknownProviderError{}
}

func getK8sNode(ctx context.Context, client client.Client) (*corev1.Node, bool, error) {
	nodes := &corev1.NodeList{}

	if err := client.List(ctx, nodes); err != nil {
		return nil, false, err
	}

	if len(nodes.Items) == 0 {
		return nil, false, nil
	}

	return &nodes.Items[0], true, nil
}

func detectNodeByProviderID(ctx context.Context, client client.Client, node *corev1.Node, scheme string) (bool, *corev1.Node, error) {
	var found bool
	var err error

	if node == nil {
		node, found, err = getK8sNode(ctx, client)
		if err != nil {
			return false, nil, err
		}
		if !found {
			return false, nil, nil
		}
	}

	if node.Spec.ProviderID != "" {
		u, err := url.Parse(node.Spec.ProviderID)
		if err != nil {
			return false, nil, err
		}

		if u.Scheme == scheme {
			return true, node, nil
		}
	}

	return false, node, nil
}
