// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package util

import (
	"net/url"
	"strings"

	"emperror.dev/errors"
	"k8s.io/client-go/tools/clientcmd"
	k8sclientapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/yaml"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
)

func GetKubeconfigWithSAToken(name, username, endpointURL string, caData []byte, saToken string) (string, error) {
	address, err := getURLWithHTTPSScheme(endpointURL)
	if err != nil {
		return "", errors.WithStack(err)
	}

	config := k8sclientapiv1.Config{
		APIVersion: k8sclientapiv1.SchemeGroupVersion.Version,
		Kind:       "Config",
		Clusters: []k8sclientapiv1.NamedCluster{
			{
				Name: name,
				Cluster: k8sclientapiv1.Cluster{
					CertificateAuthorityData: caData,
					Server:                   address,
				},
			},
		},
		Contexts: []k8sclientapiv1.NamedContext{
			{
				Name: name,
				Context: k8sclientapiv1.Context{
					Cluster:  name,
					AuthInfo: username,
				},
			},
		},
		CurrentContext: name,
		AuthInfos: []k8sclientapiv1.NamedAuthInfo{
			{
				Name: username,
				AuthInfo: k8sclientapiv1.AuthInfo{
					Token: saToken,
				},
			},
		},
	}

	y, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(y), nil
}

func GetEndpointForClusterByNetwork(cluster *clusterregistryv1alpha1.Cluster, networkName string) clusterregistryv1alpha1.KubernetesAPIEndpoint {
	var endpoint clusterregistryv1alpha1.KubernetesAPIEndpoint

	for _, apiEndpoint := range cluster.Spec.KubernetesAPIEndpoints {
		if apiEndpoint.ClientNetwork == networkName {
			endpoint = apiEndpoint

			break
		}
		// use for every network if the endpoint is not network specific
		if apiEndpoint.ClientNetwork == "" {
			endpoint = apiEndpoint
		}
	}

	return endpoint
}

func GetKubeconfigOverridesForClusterByNetwork(cluster *clusterregistryv1alpha1.Cluster, networkName string) (*clientcmd.ConfigOverrides, error) {
	overrides := &clientcmd.ConfigOverrides{}

	if len(cluster.Spec.KubernetesAPIEndpoints) == 0 {
		return overrides, nil
	}

	endpoint := GetEndpointForClusterByNetwork(cluster, networkName)

	if endpoint.ServerAddress != "" {
		address, err := getURLWithHTTPSScheme(endpoint.ServerAddress)
		if err != nil {
			return overrides, errors.WithStack(err)
		}
		overrides.ClusterInfo.Server = address
	}

	if len(endpoint.CABundle) > 0 {
		overrides.ClusterInfo.CertificateAuthorityData = endpoint.CABundle
	}

	return overrides, nil
}

func getURLWithHTTPSScheme(address string) (string, error) {
	if !strings.Contains(address, "//") {
		address = "//" + address
	}
	u, err := url.Parse(address)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	return u.String(), nil
}
