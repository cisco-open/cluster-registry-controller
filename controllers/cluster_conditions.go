// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package controllers

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
)

type ClusterConditionsMap map[clusterregistryv1alpha1.ClusterConditionType]clusterregistryv1alpha1.ClusterCondition

func SetCondition(cluster *clusterregistryv1alpha1.Cluster, currentConditions ClusterConditionsMap, condition clusterregistryv1alpha1.ClusterCondition, recorder record.EventRecorder) {
	var ok bool
	var currentCondition clusterregistryv1alpha1.ClusterCondition

	if currentCondition, ok = currentConditions[condition.Type]; !ok {
		currentCondition = condition
	} else {
		if currentCondition.Status != condition.Status {
			currentCondition.LastTransitionTime = metav1.NewTime(time.Now())
			eventType := corev1.EventTypeNormal
			if (condition.TrueIsFailure && condition.Status == corev1.ConditionTrue) || (!condition.TrueIsFailure && condition.Status == corev1.ConditionFalse) {
				eventType = corev1.EventTypeWarning
			}
			recorder.Event(cluster, eventType, condition.Reason, condition.Message)
		}
		currentCondition.Message = condition.Message
		currentCondition.Reason = condition.Reason
		currentCondition.Status = condition.Status
	}

	if currentCondition.LastTransitionTime.IsZero() {
		currentCondition.LastTransitionTime = metav1.NewTime(time.Now())
	}
	currentCondition.LastHeartbeatTime = metav1.NewTime(time.Now())

	currentConditions[condition.Type] = currentCondition
}

func GetCurrentConditions(cluster *clusterregistryv1alpha1.Cluster) ClusterConditionsMap {
	currentConditions := make(ClusterConditionsMap)
	for _, condition := range cluster.Status.Conditions {
		currentConditions[condition.Type] = condition
	}

	return currentConditions
}

func GetCurrentCondition(cluster *clusterregistryv1alpha1.Cluster, t clusterregistryv1alpha1.ClusterConditionType) clusterregistryv1alpha1.ClusterCondition {
	for _, condition := range cluster.Status.Conditions {
		if condition.Type == t {
			return condition
		}
	}

	return clusterregistryv1alpha1.ClusterCondition{}
}

func LocalClusterCondition(isClusterLocal bool) clusterregistryv1alpha1.ClusterCondition {
	condition := clusterregistryv1alpha1.ClusterCondition{
		Type:   clusterregistryv1alpha1.ClusterConditionTypeLocalCluster,
		Status: corev1.ConditionUnknown,
	}

	if isClusterLocal {
		condition.Reason = "ClusterIsLocal"
		condition.Message = "cluster is local"
		condition.Status = corev1.ConditionTrue

		return condition
	}

	condition.Reason = "ClusterIsNotLocal"
	condition.Message = "cluster is not local"
	condition.Status = corev1.ConditionFalse

	return condition
}

func LocalClusterConflictCondition(conflict bool) clusterregistryv1alpha1.ClusterCondition {
	condition := clusterregistryv1alpha1.ClusterCondition{
		Type:   clusterregistryv1alpha1.ClusterConditionTypeLocalConflict,
		Status: corev1.ConditionUnknown,

		TrueIsFailure: true,
	}

	if conflict {
		condition.Reason = "LocalClusterIsConflicting"
		condition.Message = ErrLocalClusterConflict.Error()
		condition.Status = corev1.ConditionTrue

		return condition
	}

	condition.Reason = "LocalClusterIsNotConflicting"
	condition.Message = "only one local cluster is defined"
	condition.Status = corev1.ConditionFalse

	return condition
}

func ClusterMetadataCondition(err error) clusterregistryv1alpha1.ClusterCondition {
	condition := clusterregistryv1alpha1.ClusterCondition{
		Type:   clusterregistryv1alpha1.ClusterConditionTypeClusterMetadata,
		Status: corev1.ConditionUnknown,
	}

	if err == nil {
		condition.Reason = "ClusterMetadataSet"
		condition.Message = "cluster metadata is set"
		condition.Status = corev1.ConditionTrue

		return condition
	}

	condition.Reason = "ClusterMetadataIsNotSet"
	condition.Message = err.Error()
	condition.Status = corev1.ConditionFalse

	return condition
}

func ClusterReadyCondition(err error) clusterregistryv1alpha1.ClusterCondition {
	condition := clusterregistryv1alpha1.ClusterCondition{
		Type:   clusterregistryv1alpha1.ClusterConditionTypeReady,
		Status: corev1.ConditionUnknown,
	}

	if err == nil {
		condition.Reason = "ClusterIsReady"
		condition.Message = "cluster is ready"
		condition.Status = corev1.ConditionTrue

		return condition
	}

	condition.Reason = "ClusterIsNotReady"
	condition.Message = err.Error()
	condition.Status = corev1.ConditionFalse

	return condition
}

func ClustersSyncedCondition(err error) clusterregistryv1alpha1.ClusterCondition {
	condition := clusterregistryv1alpha1.ClusterCondition{
		Type:   clusterregistryv1alpha1.ClusterConditionTypeClustersSynced,
		Status: corev1.ConditionUnknown,
	}

	if err == nil {
		condition.Reason = "ClustersAreSynced"
		condition.Message = "all participating clusters are in sync"
		condition.Status = corev1.ConditionTrue

		return condition
	}

	condition.Reason = "ClustersAreNotSynced"
	condition.Message = err.Error()
	condition.Status = corev1.ConditionFalse

	return condition
}
