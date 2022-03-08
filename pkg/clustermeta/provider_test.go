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

func TestProviderDetector(t *testing.T) {
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
	}{
		{
			filename: "amazon-eks.yaml",
			provider: "amazon",
		},
		{
			filename: "amazon-pke.yaml",
			provider: "amazon",
		},
		{
			filename: "azure-aks.yaml",
			provider: "azure",
		},
		{
			filename: "azure-pke.yaml",
			provider: "azure",
		},
		{
			filename: "gcp-gke.yaml",
			provider: "google",
		},
		{
			filename: "vsphere-pke.yaml",
			provider: "vsphere",
		},
		{
			filename: "kind-kind.yaml",
			provider: "kind",
		},
		{
			filename: "cisco-iks.yaml",
			provider: "cisco",
		},
	}

	// with mocked client
	for _, tc := range testCases {
		tc := tc
		t.Run("withclient-"+tc.filename, func(t *testing.T) {
			s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme, scheme)
			o, _, err := s.Decode(testFiles[tc.filename], nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			nodeListerClient := fake.NewFakeClientWithScheme(scheme, o)
			foundProvider, err := clustermeta.DetectProvider(context.Background(), nodeListerClient, nil)
			if err != nil {
				t.Fatal(err)
			}
			if foundProvider != tc.provider {
				t.Fatalf("%s not detected as '%s' and not '%s'", tc.filename, foundProvider, tc.provider)
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
			foundProvider, err := clustermeta.DetectProvider(context.Background(), nil, &nodes.Items[0])
			if err != nil {
				t.Fatal(err)
			}
			if foundProvider != tc.provider {
				t.Fatalf("%s not detected as '%s' and not '%s'", tc.filename, foundProvider, tc.provider)
			}
		})
	}
}

func TestUnknownProviderDetector(t *testing.T) {
	t.Parallel()
	dirName := "testdata/"
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		t.Fatal(err)
	}

	testFiles := make(map[string][]byte)
	for _, file := range files {
		testFiles[file.Name()] = ReadFile(t, dirName+file.Name())
	}

	testCases := []struct {
		filename string
		provider string
	}{{
		filename: "unknown-provider.yaml",
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
			_, err = clustermeta.DetectProvider(context.Background(), nodeListerClient, nil)
			if err == nil {
				t.Fatal("unknown provider detection ran without error")
			}
			if !clustermeta.IsUnknownProviderError(err) {
				t.Fatalf("invalid error: %v", err)
			}
		})
	}
}
