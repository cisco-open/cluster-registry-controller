// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package clustermeta_test

import (
	"context"
	"io/ioutil"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/clustermeta"
)

const testdataDir = "testdata/"

func TestDistributionDetector(t *testing.T) {
	t.Parallel()

	files, err := ioutil.ReadDir(testdataDir)
	if err != nil {
		t.Fatal(err)
	}

	testFiles := make(map[string][]byte)
	for _, file := range files {
		testFiles[file.Name()] = ReadFile(t, testdataDir+file.Name())
	}

	testCases := []struct {
		filename     string
		distribution string
	}{
		{
			filename:     "amazon-eks.yaml",
			distribution: "EKS",
		},
		{
			filename:     "amazon-pke.yaml",
			distribution: "PKE",
		},
		{
			filename:     "azure-aks.yaml",
			distribution: "AKS",
		},
		{
			filename:     "azure-pke.yaml",
			distribution: "PKE",
		},
		{
			filename:     "gcp-gke.yaml",
			distribution: "GKE",
		},
		{
			filename:     "vsphere-pke.yaml",
			distribution: "PKE",
		},
		{
			filename:     "kind-kind.yaml",
			distribution: "KIND",
		},
		{
			filename:     "cisco-iks.yaml",
			distribution: "IKS",
		},
	}

	// with mocked client
	for _, tc := range testCases {
		tc := tc
		t.Run("withclient-"+tc.filename, func(t *testing.T) {
			t.Parallel()

			s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme, scheme)
			o, _, err := s.Decode(testFiles[tc.filename], nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			nodeListerClient := fake.NewFakeClientWithScheme(scheme, o)
			foundDistribution, err := clustermeta.DetectDistribution(context.Background(), nodeListerClient, nil)
			if err != nil {
				t.Fatal(err)
			}
			if foundDistribution != tc.distribution {
				t.Fatalf("%s detected as '%s' and not '%s'", tc.filename, foundDistribution, tc.distribution)
			}
		})
	}

	// with node instane
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.filename, func(t *testing.T) {
			t.Parallel()
			s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme, scheme)
			o, _, err := s.Decode(testFiles[tc.filename], nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			nodeListerClient := fake.NewFakeClientWithScheme(scheme, o)
			nodes := &corev1.NodeList{}
			err = nodeListerClient.List(context.Background(), nodes)
			if err != nil {
				t.Fatal(err)
			}
			foundDistribution, err := clustermeta.DetectDistribution(context.Background(), nil, &nodes.Items[0])
			if err != nil {
				t.Fatal(err)
			}
			if foundDistribution != tc.distribution {
				t.Fatalf("%s detected as '%s' and not '%s'", tc.filename, foundDistribution, tc.distribution)
			}
		})
	}
}

func TestUnknownDistributionDetector(t *testing.T) {
	t.Parallel()

	files, err := ioutil.ReadDir(testdataDir)
	if err != nil {
		t.Fatal(err)
	}

	testFiles := make(map[string][]byte)
	for _, file := range files {
		testFiles[file.Name()] = ReadFile(t, testdataDir+file.Name())
	}

	testCases := []struct {
		filename string
		provider string
	}{{
		filename: "unknown-distribution.yaml",
		provider: "unknown",
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.filename, func(t *testing.T) {
			t.Parallel()

			s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme, scheme)
			o, _, err := s.Decode(testFiles[tc.filename], nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			nodeListerClient := fake.NewFakeClientWithScheme(scheme, o)
			d, err := clustermeta.DetectDistribution(context.Background(), nodeListerClient, nil)
			if err == nil {
				t.Log(d)
				t.Fatal("unknown distribution detection ran without error")
			}
			if !clustermeta.IsUnknownDistributionError(err) {
				t.Fatalf("invalid error: %v", err)
			}
		})
	}
}
