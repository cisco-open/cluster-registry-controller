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
	"bytes"
	"context"
	"net"
	"os"

	"emperror.dev/errors"

	v1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterregistryv1alpha1 "github.com/cisco-open/cluster-registry-controller/api/v1alpha1"
	"github.com/cisco-open/cluster-registry-controller/pkg/common"
)

const (
	defaultCACertPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
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

func GetK8sServerVersion(kubeConfig *rest.Config) (string, error) {
	d, err := discovery.NewDiscoveryClientForConfig(kubeConfig)
	if err != nil {
		return "", errors.WithStackIf(err)
	}

	rawVersion, err := d.ServerVersion()

	return rawVersion.String(), errors.WrapIf(err, "could not to get cluster's kubernetes apiserver version")
}

func GetReaderSecretForCluster(ctx context.Context, kubeClient client.Client, kubeConfig *rest.Config, cluster *clusterregistryv1alpha1.Cluster, saRef types.NamespacedName, apiServerEndpointAddress string) (*corev1.Secret, error) {
	sa := &corev1.ServiceAccount{}
	err := kubeClient.Get(ctx, saRef, sa)
	if err != nil {
		return nil, errors.WithStackIf(err)
	}

	// After K8s v1.24, Secret objects containing ServiceAccount tokens are no longer auto-generated, so we will have to use Token API to get the tokens.
	// Reference: https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.24.md#no-really-you-must-read-this-before-you-upgrade
	apiServerVersion, err := GetK8sServerVersion(kubeConfig)
	if err != nil {
		return nil, errors.WithStackIf(err)
	}

	isGreaterVersion, _ := common.ValidateVersionWithConstraint(apiServerVersion, ">= 1.24.0-0")

	// we will need to fetch these credentials from kubeconfig or with token api
	var saToken string
	var caData []byte

	if len(sa.Secrets) == 0 { // nolint:nestif
		if isGreaterVersion {
			k8sClient, _ := kubernetes.NewForConfig(ctrl.GetConfigOrDie())
			serviceAccounts := k8sClient.CoreV1().ServiceAccounts("cluster-registry")

			req := v1.TokenRequest{
				Spec: v1.TokenRequestSpec{
					ExpirationSeconds: pointer.Int64(60 * 60 * 24 * 365),
				},
			}

			token, err := serviceAccounts.CreateToken(context.Background(), "cluster-registry-controller-reader", &req, metav1.CreateOptions{})
			if err != nil {
				return nil, errors.WrapIf(err, "could not request token")
			}

			caCert, err := os.ReadFile(defaultCACertPath)
			if err != nil {
				return nil, errors.WrapIf(err, "could not read CA certificate file")
			}

			caData = caCert
			saToken = token.Status.Token
		} else {
			return nil, errors.NewWithDetails("could not find secret reference for sa", "sa", saRef)
		}
	} else {
		secret := &corev1.Secret{}
		err = kubeClient.Get(ctx, types.NamespacedName{
			Name:      sa.Secrets[0].Name,
			Namespace: sa.GetNamespace(),
		}, secret)
		if err != nil {
			return nil, errors.WithStackIf(err)
		}

		caData = secret.Data["ca.crt"]
		if !bytes.Contains(caData, kubeConfig.CAData) {
			caData = append(append(caData, []byte("\n")...), kubeConfig.CAData...)
		}

		saToken = string(secret.Data["token"])
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

	kubeconfig, err := GetKubeconfigWithSAToken(cluster.GetName(), sa.GetName(), apiServerEndpointAddress, caData, saToken)
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
