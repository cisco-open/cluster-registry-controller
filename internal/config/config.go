// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package config

type Configuration struct {
	MetricsAddr                string            `mapstructure:"metrics-addr" json:"metricsAddr,omitempty"`
	LeaderElection             LeaderElection    `mapstructure:"leader-election" json:"leaderElection,omitempty"`
	Logging                    Logging           `mapstructure:"log" json:"logging,omitempty"`
	ClusterController          ClusterController `mapstructure:"clusterController" json:"clusterController,omitempty"`
	SyncController             SyncController    `mapstructure:"syncController" json:"syncController,omitempty"`
	Namespace                  string            `mapstructure:"namespace" json:"namespace,omitempty"`
	ProvisionLocalCluster      string            `mapstructure:"provision-local-cluster" json:"provisionLocalCluster,omitempty"`
	ManageLocalClusterSecret   bool              `mapstructure:"manage-local-cluster-secret" json:"manageLocalClusterSecret,omitempty"`
	ReaderServiceAccountName   string            `mapstructure:"reader-service-account-name" json:"readerServiceAccountName,omitempty"`
	NetworkName                string            `mapstructure:"network-name" json:"networkName,omitempty"`
	APIServerEndpointAddress   string            `mapstructure:"apiserver-endpoint-address" json:"apiServerEndpointAddress,omitempty"`
	CoreResourcesSourceEnabled bool              `mapstructure:"core-resources-source-enabled" json:"coreResourcesSourceEnabled,omitempty"`
}

type ClusterController struct {
	WorkerCount            int `mapstructure:"workerCount" json:"workerCount,omitempty"`
	RefreshIntervalSeconds int `mapstructure:"refreshIntervalSeconds" json:"refreshIntervalSeconds,omitempty"`
}

type SyncController struct {
	WorkerCount int                     `mapstructure:"workerCount" json:"workerCount,omitempty"`
	RateLimit   SyncControllerRateLimit `mapstructure:"rateLimit" json:"rateLimit,omitempty"`
}

type SyncControllerRateLimit struct {
	MaxKeys          int `mapstructure:"maxKeys" json:"maxKeys,omitempty"`
	MaxRatePerSecond int `mapstructure:"maxRatePerSecond" json:"maxRatePerSecond,omitempty"`
	MaxBurst         int `mapstructure:"maxBurst" json:"maxBurst,omitempty"`
}

type LeaderElection struct {
	Enabled   bool   `mapstructure:"enabled" json:"enabled,omitempty"`
	Name      string `mapstructure:"name" json:"name,omitempty"`
	Namespace string `mapstructure:"namespace" json:"namespace,omitempty"`
}

type (
	LogFormat string
)

const (
	LogFormatConsole LogFormat = "console"
	LogFormatJSON    LogFormat = "json"
)

type Logging struct {
	Verbosity int8      `mapstructure:"verbosity" json:"level,omitempty"`
	Format    LogFormat `mapstructure:"format" json:"format,omitempty"`
}
