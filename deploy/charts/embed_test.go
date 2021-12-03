// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

package charts_test

import (
	"io"
	"os"
	"testing"

	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/deploy/charts"
)

func TestEmbed(t *testing.T) {
	file, err := charts.ClusterRegistry.Open("Chart.yaml")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	embeddedContent, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	localContent, err := os.ReadFile("cluster-registry/Chart.yaml")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if string(embeddedContent) != string(localContent) {
		t.Fatalf("embedded content %s does not equal local content %s", string(embeddedContent), string(localContent))
	}
}
