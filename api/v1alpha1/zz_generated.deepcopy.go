//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022 Cisco Systems, Inc. and/or its affiliates.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/banzaicloud/operator-tools/pkg/resources"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnnotationMutations) DeepCopyInto(out *AnnotationMutations) {
	*out = *in
	if in.Add != nil {
		in, out := &in.Add, &out.Add
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Remove != nil {
		in, out := &in.Remove, &out.Remove
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnnotationMutations.
func (in *AnnotationMutations) DeepCopy() *AnnotationMutations {
	if in == nil {
		return nil
	}
	out := new(AnnotationMutations)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnnotationSelector) DeepCopyInto(out *AnnotationSelector) {
	*out = *in
	if in.MatchAnnotations != nil {
		in, out := &in.MatchAnnotations, &out.MatchAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.MatchExpressions != nil {
		in, out := &in.MatchExpressions, &out.MatchExpressions
		*out = make([]AnnotationSelectorRequirement, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnnotationSelector.
func (in *AnnotationSelector) DeepCopy() *AnnotationSelector {
	if in == nil {
		return nil
	}
	out := new(AnnotationSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnnotationSelectorRequirement) DeepCopyInto(out *AnnotationSelectorRequirement) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]AnnotationValue, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnnotationSelectorRequirement.
func (in *AnnotationSelectorRequirement) DeepCopy() *AnnotationSelectorRequirement {
	if in == nil {
		return nil
	}
	out := new(AnnotationSelectorRequirement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthInfo) DeepCopyInto(out *AuthInfo) {
	*out = *in
	out.SecretRef = in.SecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthInfo.
func (in *AuthInfo) DeepCopy() *AuthInfo {
	if in == nil {
		return nil
	}
	out := new(AuthInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Cluster) DeepCopyInto(out *Cluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Cluster.
func (in *Cluster) DeepCopy() *Cluster {
	if in == nil {
		return nil
	}
	out := new(Cluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Cluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterCondition) DeepCopyInto(out *ClusterCondition) {
	*out = *in
	in.LastHeartbeatTime.DeepCopyInto(&out.LastHeartbeatTime)
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterCondition.
func (in *ClusterCondition) DeepCopy() *ClusterCondition {
	if in == nil {
		return nil
	}
	out := new(ClusterCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterFeature) DeepCopyInto(out *ClusterFeature) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterFeature.
func (in *ClusterFeature) DeepCopy() *ClusterFeature {
	if in == nil {
		return nil
	}
	out := new(ClusterFeature)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterFeature) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterFeatureList) DeepCopyInto(out *ClusterFeatureList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterFeature, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterFeatureList.
func (in *ClusterFeatureList) DeepCopy() *ClusterFeatureList {
	if in == nil {
		return nil
	}
	out := new(ClusterFeatureList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterFeatureList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterFeatureMatch) DeepCopyInto(out *ClusterFeatureMatch) {
	*out = *in
	if in.MatchLabels != nil {
		in, out := &in.MatchLabels, &out.MatchLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.MatchExpressions != nil {
		in, out := &in.MatchExpressions, &out.MatchExpressions
		*out = make([]v1.LabelSelectorRequirement, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterFeatureMatch.
func (in *ClusterFeatureMatch) DeepCopy() *ClusterFeatureMatch {
	if in == nil {
		return nil
	}
	out := new(ClusterFeatureMatch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterFeatureSpec) DeepCopyInto(out *ClusterFeatureSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterFeatureSpec.
func (in *ClusterFeatureSpec) DeepCopy() *ClusterFeatureSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterFeatureSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterFeatureStatus) DeepCopyInto(out *ClusterFeatureStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterFeatureStatus.
func (in *ClusterFeatureStatus) DeepCopy() *ClusterFeatureStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterFeatureStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterList) DeepCopyInto(out *ClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Cluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterList.
func (in *ClusterList) DeepCopy() *ClusterList {
	if in == nil {
		return nil
	}
	out := new(ClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterMetadata) DeepCopyInto(out *ClusterMetadata) {
	*out = *in
	if in.KubeProxyVersions != nil {
		in, out := &in.KubeProxyVersions, &out.KubeProxyVersions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.KubeletVersions != nil {
		in, out := &in.KubeletVersions, &out.KubeletVersions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Locality != nil {
		in, out := &in.Locality, &out.Locality
		*out = new(Locality)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterMetadata.
func (in *ClusterMetadata) DeepCopy() *ClusterMetadata {
	if in == nil {
		return nil
	}
	out := new(ClusterMetadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterSpec) DeepCopyInto(out *ClusterSpec) {
	*out = *in
	out.AuthInfo = in.AuthInfo
	if in.KubernetesAPIEndpoints != nil {
		in, out := &in.KubernetesAPIEndpoints, &out.KubernetesAPIEndpoints
		*out = make([]KubernetesAPIEndpoint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterSpec.
func (in *ClusterSpec) DeepCopy() *ClusterSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterStatus) DeepCopyInto(out *ClusterStatus) {
	*out = *in
	in.ClusterMetadata.DeepCopyInto(&out.ClusterMetadata)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]ClusterCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterStatus.
func (in *ClusterStatus) DeepCopy() *ClusterStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContentSelector) DeepCopyInto(out *ContentSelector) {
	*out = *in
	out.Value = in.Value
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContentSelector.
func (in *ContentSelector) DeepCopy() *ContentSelector {
	if in == nil {
		return nil
	}
	out := new(ContentSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubernetesAPIEndpoint) DeepCopyInto(out *KubernetesAPIEndpoint) {
	*out = *in
	if in.CABundle != nil {
		in, out := &in.CABundle, &out.CABundle
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesAPIEndpoint.
func (in *KubernetesAPIEndpoint) DeepCopy() *KubernetesAPIEndpoint {
	if in == nil {
		return nil
	}
	out := new(KubernetesAPIEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LabelMutations) DeepCopyInto(out *LabelMutations) {
	*out = *in
	if in.Add != nil {
		in, out := &in.Add, &out.Add
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Remove != nil {
		in, out := &in.Remove, &out.Remove
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LabelMutations.
func (in *LabelMutations) DeepCopy() *LabelMutations {
	if in == nil {
		return nil
	}
	out := new(LabelMutations)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Locality) DeepCopyInto(out *Locality) {
	*out = *in
	if in.Regions != nil {
		in, out := &in.Regions, &out.Regions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Zones != nil {
		in, out := &in.Zones, &out.Zones
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Locality.
func (in *Locality) DeepCopy() *Locality {
	if in == nil {
		return nil
	}
	out := new(Locality)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in MatchedRules) DeepCopyInto(out *MatchedRules) {
	{
		in := &in
		*out = make(MatchedRules, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchedRules.
func (in MatchedRules) DeepCopy() MatchedRules {
	if in == nil {
		return nil
	}
	out := new(MatchedRules)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mutations) DeepCopyInto(out *Mutations) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = new(AnnotationMutations)
		(*in).DeepCopyInto(*out)
	}
	if in.GVK != nil {
		in, out := &in.GVK, &out.GVK
		*out = new(resources.GroupVersionKind)
		**out = **in
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = new(LabelMutations)
		(*in).DeepCopyInto(*out)
	}
	if in.Overrides != nil {
		in, out := &in.Overrides, &out.Overrides
		*out = make([]resources.K8SResourceOverlayPatch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mutations.
func (in *Mutations) DeepCopy() *Mutations {
	if in == nil {
		return nil
	}
	out := new(Mutations)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacedName) DeepCopyInto(out *NamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacedName.
func (in *NamespacedName) DeepCopy() *NamespacedName {
	if in == nil {
		return nil
	}
	out := new(NamespacedName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSyncRule) DeepCopyInto(out *ResourceSyncRule) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSyncRule.
func (in *ResourceSyncRule) DeepCopy() *ResourceSyncRule {
	if in == nil {
		return nil
	}
	out := new(ResourceSyncRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceSyncRule) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSyncRuleList) DeepCopyInto(out *ResourceSyncRuleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceSyncRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSyncRuleList.
func (in *ResourceSyncRuleList) DeepCopy() *ResourceSyncRuleList {
	if in == nil {
		return nil
	}
	out := new(ResourceSyncRuleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceSyncRuleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSyncRuleSpec) DeepCopyInto(out *ResourceSyncRuleSpec) {
	*out = *in
	if in.ClusterFeatureMatches != nil {
		in, out := &in.ClusterFeatureMatches, &out.ClusterFeatureMatches
		*out = make([]ClusterFeatureMatch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.GVK = in.GVK
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]SyncRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSyncRuleSpec.
func (in *ResourceSyncRuleSpec) DeepCopy() *ResourceSyncRuleSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceSyncRuleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSyncRuleStatus) DeepCopyInto(out *ResourceSyncRuleStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSyncRuleStatus.
func (in *ResourceSyncRuleStatus) DeepCopy() *ResourceSyncRuleStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceSyncRuleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SyncRule) DeepCopyInto(out *SyncRule) {
	*out = *in
	if in.Matches != nil {
		in, out := &in.Matches, &out.Matches
		*out = make([]SyncRuleMatch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Mutations.DeepCopyInto(&out.Mutations)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SyncRule.
func (in *SyncRule) DeepCopy() *SyncRule {
	if in == nil {
		return nil
	}
	out := new(SyncRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SyncRuleMatch) DeepCopyInto(out *SyncRuleMatch) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make([]AnnotationSelector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Content != nil {
		in, out := &in.Content, &out.Content
		*out = make([]ContentSelector, len(*in))
		copy(*out, *in)
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]v1.LabelSelector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.ObjectKey = in.ObjectKey
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SyncRuleMatch.
func (in *SyncRuleMatch) DeepCopy() *SyncRuleMatch {
	if in == nil {
		return nil
	}
	out := new(SyncRuleMatch)
	in.DeepCopyInto(out)
	return out
}
