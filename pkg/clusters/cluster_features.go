// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package clusters

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type ClusterFeature interface {
	GetUID() string
	GetName() string
	GetLabels() map[string]string
}

type clusterFeature struct {
	UID    string
	Name   string
	Labels map[string]string
}

func (f clusterFeature) GetUID() string {
	return f.UID
}

func (f clusterFeature) GetName() string {
	return f.Name
}

func (f clusterFeature) GetLabels() map[string]string {
	return f.Labels
}

func NewClusterFeature(uid, name string, labels map[string]string) ClusterFeature {
	return clusterFeature{
		UID:    uid,
		Name:   name,
		Labels: labels,
	}
}

type ClusterFeatureRequirement struct {
	Name             string
	MatchLabels      map[string]string
	MatchExpressions []metav1.LabelSelectorRequirement
}

func (r ClusterFeatureRequirement) Match(features map[string]ClusterFeature) bool {
	for _, feature := range features {
		if r.Name != "" && r.Name != feature.GetName() {
			continue
		}

		if len(r.MatchLabels) == 0 && len(r.MatchExpressions) == 0 {
			return true
		}

		matcher, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
			MatchLabels:      r.MatchLabels,
			MatchExpressions: r.MatchExpressions,
		})
		if err != nil {
			return false
		}
		if matcher.Matches(labels.Set(feature.GetLabels())) {
			return true
		}
	}

	return false
}
