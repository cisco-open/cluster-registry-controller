// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package charts_test

import (
	"io"
	"os"
	"testing"

	"github.com/cisco-open/cluster-registry-controller/deploy/charts"
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
