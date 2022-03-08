// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package util

import (
	"bytes"
	"context"
	"net"

	"emperror.dev/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
)

func GetExternalAddressOfAPIServer(kubeConfig *rest.Config) (string, error) {
	d, err := discovery.NewDiscoveryClientForConfig(kubeConfig)
	if err != nil {
		return "", errors.WithStackIf(err)
	}

	v := &metav1.APIVersions{}
	err = d.RESTClient().Get().AbsPath(d.LegacyPrefix).Do(context.TODO()).Into(v)
	if err != nil {
		return "", errors.WithStackIf(err)
	}

	for _, addr := range v.ServerAddressByClientCIDRs {
		if addr.ClientCIDR == (&net.IPNet{
			IP:   net.IPv4zero,
			Mask: net.IPv4Mask(0, 0, 0, 0),
		}).String() {
			return addr.ServerAddress, nil
		}
	}

	return "", errors.New("could not determine external apiserver address")
}

func GetReaderSecretForCluster(ctx context.Context, kubeClient client.Client, kubeConfig *rest.Config, cluster *clusterregistryv1alpha1.Cluster, saRef types.NamespacedName, apiServerEndpointAddress string) (*corev1.Secret, error) {
	sa := &corev1.ServiceAccount{}
	err := kubeClient.Get(ctx, saRef, sa)
	if err != nil {
		return nil, errors.WithStackIf(err)
	}

	if len(sa.Secrets) == 0 {
		return nil, errors.NewWithDetails("could not find secret reference for sa", "sa", saRef)
	}

	secret := &corev1.Secret{}
	err = kubeClient.Get(ctx, types.NamespacedName{
		Name:      sa.Secrets[0].Name,
		Namespace: sa.GetNamespace(),
	}, secret)
	if err != nil {
		return nil, errors.WithStackIf(err)
	}

	caData := secret.Data["ca.crt"]
	if !bytes.Contains(caData, kubeConfig.CAData) {
		caData = append(append(caData, []byte("\n")...), kubeConfig.CAData...)
	}

	// add overrides specified in the cluster resource without network specified
	endpoint := GetEndpointForClusterByNetwork(cluster, "")
	if endpoint.ServerAddress != "" {
		apiServerEndpointAddress = endpoint.ServerAddress
	}
	if len(endpoint.CABundle) > 0 {
		caData = append(append(caData, []byte("\n")...), endpoint.CABundle...)
	}

	// try to get external ip address from api server
	if apiServerEndpointAddress == "" {
		apiServerEndpointAddress, _ = GetExternalAddressOfAPIServer(kubeConfig)
	}

	// try to get endpoint from used kubeconfig
	if apiServerEndpointAddress == "" {
		apiServerEndpointAddress = kubeConfig.Host
	}

	kubeconfig, err := GetKubeconfigWithSAToken(cluster.GetName(), sa.GetName(), apiServerEndpointAddress, caData, string(secret.Data["token"]))
	if err != nil {
		return nil, errors.WithStackIf(err)
	}

	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Spec.AuthInfo.SecretRef.Name,
			Namespace: cluster.Spec.AuthInfo.SecretRef.Namespace,
		},
		Type: clusterregistryv1alpha1.SecretTypeClusterRegistry,
		Data: map[string][]byte{
			"kubeconfig": []byte(kubeconfig),
		},
	}, nil
}
