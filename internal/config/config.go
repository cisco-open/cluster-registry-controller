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

	// ClusterValidatorWebhook configures the cluster CR validator webhook for
	// the operator.
	ClusterValidatorWebhook ClusterValidatorWebhook `mapstructure:"cluster-validator-webhook" json:"clusterValidatorWebhook"`
}

// ClusterValidatorWebhook describes the configuration options for the cluster
// CR validator webhook.
type ClusterValidatorWebhook struct {
	// Enabled is the indicator to determine whether the webhook is enabled.
	Enabled bool `mapstructure:"enabled" json:"enabled,omitempty"`

	// Certificate directory is the directory where the cluster CR validator webhook stores its certificates locally.
	CertificateDirectory string `mapstructure:"certificate-directory" json:"certificateDirectory,omitempty"`

	// Name is the name of the cluster CR validator webhook resource.
	Name string `mapstructure:"name" json:"name,omitempty"`

	// Port is the port the cluster CR validator webhook serves on.
	Port uint `mapstructure:"port" json:"port,omitempty"`
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
