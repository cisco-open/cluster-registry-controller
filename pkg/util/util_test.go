// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package util_test

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
	"github.com/banzaicloud/operator-tools/pkg/resources"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/util"
)

func TestK8SResourceOverlayPatchExecuteTemplate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		patch   resources.K8SResourceOverlayPatch
		object  runtime.Object
		cluster *clusterregistryv1alpha1.Cluster
		wanted  string
	}{
		"working template": {
			patch: resources.K8SResourceOverlayPatch{
				Path:  utils.StringPointer("/name"),
				Type:  resources.ReplaceOverlayPatchType,
				Value: utils.StringPointer(`{{ printf "%s-%s" .Object.GetName .Cluster.GetName }}`),
			},
			object: &corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-service",
				},
			},
			cluster: &clusterregistryv1alpha1.Cluster{
				ObjectMeta: v1.ObjectMeta{
					Name: "demo-cluster-1",
				},
			},
			wanted: "test-service-demo-cluster-1",
		},
	}

	for _, test := range tests {
		result, err := util.K8SResourceOverlayPatchExecuteTemplate(test.patch, map[string]interface{}{
			"Object":  test.object,
			"Cluster": test.cluster,
		})
		if err != nil {
			t.Fatal(err)
		}
		if utils.PointerToString(result.Value) != test.wanted {
			t.Fatal(fmt.Errorf("%s != %s", utils.PointerToString(result.Value), test.wanted)) // nolint:goerr113
		}
	}
}
