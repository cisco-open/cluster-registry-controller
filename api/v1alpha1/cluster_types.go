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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	SecretTypeClusterRegistry corev1.SecretType = "k8s.cisco.com/cluster-registry-secret"
	KubeconfigKey                               = "kubeconfig"
)

// AuthInfo holds information that describes how a client can get
// credentials to access the cluster.
type AuthInfo struct {
	SecretRef NamespacedName `json:"secretRef,omitempty"`
}

// Equivalent of types.NamespacedName with JSON tags
type NamespacedName struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type KubernetesAPIEndpoint struct {
	// The network name of the client to match whether if it should
	// use the corresponding server address.
	// +optional
	ClientNetwork string `json:"clientNetwork,omitempty"`
	// Address of this server, suitable for a client that matches the clientNetwork if specified.
	// This can be a hostname, hostname:port, IP or IP:port.
	ServerAddress string `json:"serverAddress,omitempty"`
	// CABundle contains the certificate authority information.
	// +optional
	CABundle []byte `json:"caBundle,omitempty"`
}

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// UID of the kube-system namespace
	ClusterID types.UID `json:"clusterID"`
	// AuthInfo holds information that describes how a client can get
	// credentials to access the cluster.
	AuthInfo AuthInfo `json:"authInfo,omitempty"`
	// KubernetesAPIEndpoints represents the endpoints of the API server for this
	// cluster.
	// +optional
	KubernetesAPIEndpoints []KubernetesAPIEndpoint `json:"kubernetesApiEndpoints,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	State   ClusterState `json:"state,omitempty"`
	Message string       `json:"message,omitempty"`

	Type   ClusterType `json:"type,omitempty"`
	Leader bool        `json:"leader,omitempty"`

	// Metadata
	ClusterMetadata `json:",inline"`

	// Conditions contains the different condition statuses for this cluster.
	Conditions []ClusterCondition `json:"conditions,omitempty"`
}

func (s ClusterStatus) Reset() ClusterStatus {
	return ClusterStatus{
		State:           ClusterStateReady,
		Type:            ClusterTypePeer,
		ClusterMetadata: ClusterMetadata{},
		Conditions:      s.Conditions,
	}
}

type ClusterType string

const (
	ClusterTypeLocal ClusterType = "Local"
	ClusterTypePeer  ClusterType = "Peer"
)

type ClusterState string

const (
	ClusterStateReady           ClusterState = "Ready"
	ClusterStateFailed          ClusterState = "Failed"
	ClusterStateDisabled        ClusterState = "Disabled"
	ClusterStateMissingAuthInfo ClusterState = "MissingAuthInfo"
	ClusterStateInvalidAuthInfo ClusterState = "InvalidAuthInfo"
)

// ClusterConditionType marks the kind of cluster condition being reported.
type ClusterConditionType string

const (
	ClusterConditionTypeLocalCluster    ClusterConditionType = "LocalCluster"
	ClusterConditionTypeLocalConflict   ClusterConditionType = "LocalConflict"
	ClusterConditionTypeClusterMetadata ClusterConditionType = "ClusterMetadataSet"
	ClusterConditionTypeReady           ClusterConditionType = "Ready"
	ClusterConditionTypeClustersSynced  ClusterConditionType = "ClustersSynced"
)

// ClusterCondition contains condition information for a cluster.
type ClusterCondition struct {
	// Type is the type of the cluster condition.
	Type ClusterConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ClusterConditionType"`

	// Status is the status of the condition. One of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`

	// LastHeartbeatTime is the last time this condition was updated.
	// +optional
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty" protobuf:"bytes,3,opt,name=lastHeartbeatTime"`

	// LastTransitionTime is the last time the condition changed from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`

	// Reason is a (brief) reason for the condition's last status change.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`

	// Message is a human-readable message indicating details about the last status change.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`

	// Internally used flag to mark whether a true status should be handled as a failure
	TrueIsFailure bool `json:"-"`
}

// +kubebuilder:object:root=true

// Cluster is the Schema for the clusters API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=clusters,scope=Cluster
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".spec.clusterID",priority=1
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.state"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".status.type"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type==\"ClustersSynced\")].status"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".status.version",priority=1
// +kubebuilder:printcolumn:name="Provider",type="string",JSONPath=".status.provider",priority=1
// +kubebuilder:printcolumn:name="Distribution",type="string",JSONPath=".status.distribution",priority=1
// +kubebuilder:printcolumn:name="Region",type="string",JSONPath=".status.locality.region",priority=1
// +kubebuilder:printcolumn:name="Status Message",type="string",JSONPath=".status.message",priority=1
// +kubebuilder:printcolumn:name="Sync Message",type="string",JSONPath=".status.conditions[?(@.type==\"ClustersSynced\")].message",priority=1
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
