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
	"encoding/json"

	"github.com/tidwall/gjson"
	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/banzaicloud/operator-tools/pkg/resources"
)

func (m SyncRuleMatch) Match(obj runtime.Object) (bool, error) {
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		return false, err
	}

	// objectkey
	if m.ObjectKey.Name != "" && m.ObjectKey.Name != objMeta.GetName() {
		return false, nil
	}
	if m.ObjectKey.Namespace != "" && m.ObjectKey.Namespace != objMeta.GetNamespace() {
		return false, nil
	}

	// namespace
	if len(m.Namespaces) > 0 && objMeta.GetNamespace() != "" {
		found := false
		for _, name := range m.Namespaces {
			if objMeta.GetNamespace() == name {
				found = true
			}
		}
		if !found {
			return false, nil
		}
	}

	// labels
	if len(m.Labels) > 0 {
		found := false
		for _, labelSelector := range m.Labels {
			labelSelector := labelSelector
			matcher, err := metav1.LabelSelectorAsSelector(&labelSelector)
			if err != nil {
				return false, err
			}
			if matcher.Matches(labels.Set(objMeta.GetLabels())) {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}

	// annotations
	if len(m.Annotations) > 0 {
		found := false
		for _, annotationSelector := range m.Annotations {
			labelSelectorFromAnnotations := metav1.LabelSelector{
				MatchExpressions: annotationSelector.convertMatchExpressions(),
				MatchLabels:      annotationSelector.MatchAnnotations,
			}

			matcher, err := metav1.LabelSelectorAsSelector(&labelSelectorFromAnnotations)
			if err != nil {
				return false, err
			}
			if matcher.Matches(labels.Set(objMeta.GetAnnotations())) {
				found = true
			}
		}
		if !found {
			return false, nil
		}
	}

	// content
	if len(m.Content) > 0 {
		j, err := json.Marshal(obj)
		if err != nil {
			return false, err
		}

		for _, content := range m.Content {
			res := gjson.GetBytes(j, content.Key)
			if !res.Exists() {
				return false, nil
			}
			switch content.Value.Type {
			case intstr.Int:
				if res.Int() != int64(content.Value.IntVal) {
					return false, nil
				}
			case intstr.String:
				if res.String() != content.Value.StrVal {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

type MatchedRules []SyncRule

func (r ResourceSyncRuleSpec) Match(obj runtime.Object) (bool, MatchedRules, error) {
	matchedRules := make(MatchedRules, 0)

	if resources.ConvertGVK(obj.GetObjectKind().GroupVersionKind()) != r.GVK {
		return false, matchedRules, nil
	}

	for _, rule := range r.Rules {
		ok, err := rule.Match(obj)
		if err != nil {
			return false, matchedRules, err
		}
		if ok {
			matchedRules = append(matchedRules, rule)
		}
	}

	return len(matchedRules) > 0, matchedRules, nil
}

func (s *ResourceSyncRule) Match(obj runtime.Object) (bool, MatchedRules, error) {
	return s.Spec.Match(obj)
}

func (r *SyncRule) Match(obj runtime.Object) (bool, error) {
	if len(r.Matches) == 0 {
		return true, nil
	}

	for _, m := range r.Matches {
		if ok, err := m.Match(obj); ok && err == nil {
			return ok, nil
		} else if err != nil {
			return false, err
		}
	}

	return false, nil
}

func (r MatchedRules) GetMutatedGVK(gvk schema.GroupVersionKind) (bool, schema.GroupVersionKind) {
	var mutated bool

	for _, matchedRule := range r {
		mGVK := schema.GroupVersionKind(matchedRule.Mutations.GetGVK())
		if !mGVK.Empty() {
			mutated = true
			if mGVK.Group != "" {
				gvk.Group = mGVK.Group
			}
			if mGVK.Kind != "" {
				gvk.Kind = mGVK.Kind
			}
			if mGVK.Version != "" {
				gvk.Version = mGVK.Version
			}
		}
	}

	return mutated, gvk
}

func (r MatchedRules) GetMutationLabels() LabelMutations {
	m := LabelMutations{
		Add:    make(map[string]string),
		Remove: make([]string, 0),
	}

	for _, matchedRule := range r {
		for k, v := range matchedRule.Mutations.GetLabels().Add {
			m.Add[k] = v
		}
		m.Remove = append(m.Remove, matchedRule.Mutations.GetLabels().Remove...)
	}

	return m
}

func (r MatchedRules) GetMutationAnnotations() AnnotationMutations {
	m := AnnotationMutations{
		Add:    make(map[string]string),
		Remove: make([]string, 0),
	}

	for _, matchedRule := range r {
		for k, v := range matchedRule.Mutations.GetAnnotations().Add {
			m.Add[k] = v
		}
		m.Remove = append(m.Remove, matchedRule.Mutations.GetAnnotations().Remove...)
	}

	return m
}

func (r MatchedRules) GetMutationOverrides() []resources.K8SResourceOverlayPatch {
	overrides := make([]resources.K8SResourceOverlayPatch, 0)
	for _, matchedRule := range r {
		overrides = append(overrides, matchedRule.Mutations.Overrides...)
	}

	return overrides
}

func (r MatchedRules) GetMutationSyncStatus() bool {
	for _, matchedRule := range r {
		if matchedRule.Mutations.SyncStatus == true {
			return true
		}
	}

	return false
}

func (s AnnotationSelector) convertMatchExpressions() []metav1.LabelSelectorRequirement {
	reqs := make([]metav1.LabelSelectorRequirement, 0)
	for _, r := range s.MatchExpressions {
		values := make([]string, len(r.Values))
		for i, v := range r.Values {
			values[i] = string(v)
		}
		reqs = append(reqs, metav1.LabelSelectorRequirement{
			Key:      r.Key,
			Operator: r.Operator,
			Values:   values,
		})
	}

	return reqs
}
