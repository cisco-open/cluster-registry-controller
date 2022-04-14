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

package clusters_test

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/cisco-open/cluster-registry-controller/pkg/clusters"
)

func TestClusterFeatures(t *testing.T) {
	t.Parallel()

	cl, err := clusters.NewCluster(context.Background(), "test", &rest.Config{}, logr.Discard())
	if err != nil {
		panic(err)
	}

	cf := clusters.NewClusterFeature("test-feature", "test-feature",
		map[string]string{
			"testlabel":    "testlabelvalue",
			"testlabelkey": "something",
		},
	)

	cf2 := clusters.NewClusterFeature("another-test-feature", "another-test-feature", nil)

	r := clusters.NewManagedReconciler("test", logr.Discard())
	c := clusters.NewManagedController("test", r, logr.Discard(), clusters.WithRequiredClusterFeatures(
		clusters.ClusterFeatureRequirement{
			Name: "test-feature",
			MatchLabels: map[string]string{
				"testlabel": "testlabelvalue",
			},
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "testlabelkey",
					Operator: metav1.LabelSelectorOpExists,
				},
			},
		},
		clusters.ClusterFeatureRequirement{
			Name: "another-test-feature",
		},
	))

	err = cl.AddController(c)
	if err != nil {
		panic(err)
	}

	if len(cl.GetControllers()) != 0 {
		t.Fatalf("controllers count != 0")
	}

	if len(cl.GetPendingControllers()) != 1 {
		t.Fatalf("pending controllers count != 1")
	}

	cl.AddFeature(cf)

	if len(cl.GetControllers()) != 0 {
		t.Fatalf("controllers count != 0 after the first feature is added")
	}
	if len(cl.GetPendingControllers()) != 1 {
		t.Fatalf("pending controllers != 1 after the first feature is added")
	}

	cl.AddFeature(cf2)
	if len(cl.GetControllers()) != 1 {
		t.Fatalf("controllers count != 1 after the second feature is added")
	}
	if len(cl.GetPendingControllers()) != 0 {
		t.Fatalf("pending controllers != 0 after the second feature is added")
	}

	cl.RemoveFeature(cf.GetUID())
	if len(cl.GetControllers()) != 0 {
		t.Fatalf("controllers count != 0 after the feature is removed")
	}
	if len(cl.GetPendingControllers()) != 1 {
		t.Fatalf("pending controllers != 1 after the feature is removed")
	}

	cl.RemoveController(c)
	if len(cl.GetControllers()) != 0 {
		t.Fatalf("controllers count != 0 after the controller is removed")
	}
	if len(cl.GetPendingControllers()) != 0 {
		t.Fatalf("pending controllers != 0 after the controller is removed")
	}
}
