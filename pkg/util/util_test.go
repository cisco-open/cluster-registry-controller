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

package util_test

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/banzaicloud/operator-tools/pkg/resources"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	clusterregistryv1alpha1 "github.com/cisco-open/cluster-registry-controller/api/v1alpha1"
	"github.com/cisco-open/cluster-registry-controller/pkg/util"
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
