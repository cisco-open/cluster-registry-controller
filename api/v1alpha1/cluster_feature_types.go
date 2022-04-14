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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterFeatureSpec defines the desired state of ClusterFeature
type ClusterFeatureSpec struct {
	FeatureName string `json:"featureName"`
}

// ClusterFeatureStatus defines the observed state of ClusterFeature
type ClusterFeatureStatus struct{}

// +kubebuilder:object:root=true

// ClusterFeature is the Schema for the clusterfeatures API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=clusterfeatures,scope=Cluster,shortName=cf
// +kubebuilder:printcolumn:name="Feature",type="string",JSONPath=".spec.featureName"
type ClusterFeature struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterFeatureSpec   `json:"spec"`
	Status ClusterFeatureStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterFeatureList contains a list of ClusterFeature
type ClusterFeatureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterFeature `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterFeature{}, &ClusterFeatureList{})
}
