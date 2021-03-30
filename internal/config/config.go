// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

package config

type Configuration struct {
	MetricsAddr       string            `mapstructure:"metrics-addr" json:"metricsAddr,omitempty"`
	LeaderElection    LeaderElection    `mapstructure:"leader-election" json:"leaderElection,omitempty"`
	Logging           Logging           `mapstructure:"log" json:"logging,omitempty"`
	ClusterController ClusterController `mapstructure:"clusterController" json:"clusterController,omitempty"`
	SyncController    SyncController    `mapstructure:"syncController" json:"syncController,omitempty"`
	Namespace         string            `mapstructure:"namespace" json:"namespace,omitempty"`
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
