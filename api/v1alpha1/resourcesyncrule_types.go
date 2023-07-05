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
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/banzaicloud/operator-tools/pkg/resources"
	"github.com/banzaicloud/operator-tools/pkg/types"
)

const (
	OwnershipAnnotation       = "cluster-registry.k8s.cisco.com/resource-owner-cluster-id"
	OriginalGVKAnnotation     = "cluster-registry.k8s.cisco.com/original-group-version-kind"
	ClusterDisabledAnnotation = "cluster-registry.k8s.cisco.com/cluster-disabled"
	SyncDisabledAnnotation    = "cluster-registry.k8s.cisco.com/resource-sync-disabled"
)

type ResourceSyncRuleSpec struct {
	ClusterFeatureMatches []ClusterFeatureMatch      `json:"clusterFeatureMatch,omitempty"`
	GVK                   resources.GroupVersionKind `json:"groupVersionKind"`
	Rules                 []SyncRule                 `json:"rules"`
}

type ClusterFeatureMatch struct {
	FeatureName      string                            `json:"featureName,omitempty"`
	MatchLabels      map[string]string                 `json:"matchLabels,omitempty"`
	MatchExpressions []metav1.LabelSelectorRequirement `json:"matchExpressions,omitempty"`
}

type SyncRule struct {
	Matches   []SyncRuleMatch `json:"match,omitempty"`
	Mutations Mutations       `json:"mutations,omitempty"`
}

type Mutations struct {
	Annotations *AnnotationMutations                `json:"annotations,omitempty"`
	GVK         *resources.GroupVersionKind         `json:"groupVersionKind,omitempty"`
	Labels      *LabelMutations                     `json:"labels,omitempty"`
	Overrides   []resources.K8SResourceOverlayPatch `json:"overrides,omitempty"`
	SyncStatus  bool                                `json:"syncStatus,omitempty"`
}

func (m Mutations) GetGVK() resources.GroupVersionKind {
	if m.GVK != nil {
		return *m.GVK
	}

	return resources.GroupVersionKind{}
}

func (m Mutations) GetAnnotations() AnnotationMutations {
	if m.Annotations != nil {
		return *m.Annotations
	}

	return AnnotationMutations{}
}

func (m Mutations) GetLabels() LabelMutations {
	if m.Labels != nil {
		return *m.Labels
	}

	return LabelMutations{}
}

type AnnotationMutations struct {
	Add    map[string]string `json:"add,omitempty"`
	Remove []string          `json:"remove,omitempty"`
}

type LabelMutations struct {
	Add    map[string]string `json:"add,omitempty"`
	Remove []string          `json:"remove,omitempty"`
}

type SyncRuleMatch struct {
	Annotations []AnnotationSelector   `json:"annotations,omitempty"`
	Content     []ContentSelector      `json:"content,omitempty"`
	Labels      []metav1.LabelSelector `json:"labels,omitempty"`
	Namespaces  []string               `json:"namespaces,omitempty"`
	ObjectKey   types.ObjectKey        `json:"objectKey,omitempty"`
}

type AnnotationSelector struct {
	MatchAnnotations map[string]string               `json:"matchAnnotations,omitempty"`
	MatchExpressions []AnnotationSelectorRequirement `json:"matchExpressions,omitempty"`
}

type ContentSelector struct {
	Key   string             `json:"key"`
	Value intstr.IntOrString `json:"value"`
}

// +kubebuilder:validation:Pattern=`^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$`
type AnnotationValue string

// A annotation selector requirement is a selector that contains values, a key, and an operator that
// relates the key and values.
type AnnotationSelectorRequirement struct {
	// key is the label key that the selector applies to.
	// +patchMergeKey=key
	// +patchStrategy=merge
	Key string `json:"key" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,1,opt,name=key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator metav1.LabelSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=LabelSelectorOperator"`
	// values is an array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. This array is replaced during a strategic
	// merge patch.
	// +optional
	Values []AnnotationValue `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

// An annotation selector operator is the set of operators that can be used in a selector requirement.
type AnnotationSelectorOperator string

const (
	AnnotationSelectorOpExists       AnnotationSelectorOperator = "Exists"
	AnnotationSelectorOpDoesNotExist AnnotationSelectorOperator = "DoesNotExist"
)

type ResourceSyncRuleStatus struct{}

// GetRelatedNamespaces gathers namespaces from rule matchers
// it returns an empty list if any of the matchers is without namespace filter
// since it means every namespace is related
func (s *ResourceSyncRule) GetRelatedNamespaces() []string {
	namespaces := []string{}

	_namespaces := map[string]struct{}{}

	for _, rule := range s.Spec.Rules {
		for _, match := range rule.Matches {
			if len(match.Namespaces) == 0 {
				return namespaces
			}
			for _, ns := range match.Namespaces {
				_namespaces[ns] = struct{}{}
			}
		}
	}

	for ns := range _namespaces {
		namespaces = append(namespaces, ns)
	}

	return namespaces
}

// +kubebuilder:object:root=true

// ResourceSyncRule is the Schema for the resource sync rule API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=resourcesyncrules,scope=Cluster,shortName=rsr
type ResourceSyncRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceSyncRuleSpec   `json:"spec,omitempty"`
	Status ResourceSyncRuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ResourceSyncRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceSyncRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceSyncRule{}, &ResourceSyncRuleList{})
}
