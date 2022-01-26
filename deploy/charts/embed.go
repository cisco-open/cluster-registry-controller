// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

package charts

import (
	"embed"
	"io/fs"
)

var (
	//go:embed cluster-registry cluster-registry/templates/_helpers.tpl
	clusterRegistryEmbed embed.FS

	// ClusterRegistry exposes the cluster-registry chart using relative file paths from the chart root
	ClusterRegistry fs.FS
)

func init() {
	var err error
	ClusterRegistry, err = fs.Sub(clusterRegistryEmbed, "cluster-registry")
	if err != nil {
		panic(err)
	}
}