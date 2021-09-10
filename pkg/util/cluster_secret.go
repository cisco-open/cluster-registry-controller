// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

package util

import (
	"bytes"
	"context"
	"net"
	"net/url"

	"emperror.dev/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	k8sclientapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

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
			return (&url.URL{
				Scheme: "https",
				Host:   addr.ServerAddress,
			}).String(), nil
		}
	}

	return "", errors.New("could not determine external apiserver address")
}

func GetReaderSecretForCluster(ctx context.Context, kubeClient client.Client, kubeConfig *rest.Config, cluster *clusterregistryv1alpha1.Cluster, saRef types.NamespacedName) (*corev1.Secret, error) {
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

	// try to get external ip address from api server
	apiServerAddress, _ := GetExternalAddressOfAPIServer(kubeConfig)

	if apiServerAddress == "" {
		apiServerAddress = kubeConfig.Host
	}

	kubeconfig, err := GetKubeconfigWithSAToken(cluster.GetName(), sa.GetName(), apiServerAddress, caData, string(secret.Data["token"]))
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

func GetKubeconfigWithSAToken(name, username, url string, caData []byte, saToken string) (string, error) {
	config := k8sclientapiv1.Config{
		APIVersion: k8sclientapiv1.SchemeGroupVersion.Version,
		Kind:       "Config",
		Clusters: []k8sclientapiv1.NamedCluster{
			{
				Name: name,
				Cluster: k8sclientapiv1.Cluster{
					CertificateAuthorityData: caData,
					Server:                   url,
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
