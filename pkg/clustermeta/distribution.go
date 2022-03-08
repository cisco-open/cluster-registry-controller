// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package clustermeta

import (
	"context"
	"errors"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var knownDistributions = []IsDistribution{
	IsEKS,
	IsPKE,
	IsGKE,
	IsAKS,
	IsKIND,
	IsIKS,
}

type IsDistribution func(ctx context.Context, client client.Client, node *corev1.Node) (bool, string, error)

type UnknownDistributionError struct{}

func (UnknownDistributionError) Error() string {
	return "unknown distribution"
}

func IsUnknownDistributionError(err error) bool {
	return errors.As(err, &UnknownDistributionError{})
}

func DetectDistribution(ctx context.Context, client client.Client, node *corev1.Node) (string, error) {
	for _, f := range knownDistributions {
		select {
		case <-ctx.Done():
			return "", UnknownDistributionError{}
		default:
			if ok, distributionName, err := f(ctx, client, node); err != nil {
				return "", err
			} else if ok {
				return distributionName, nil
			}
		}
	}

	return "", UnknownDistributionError{}
}
