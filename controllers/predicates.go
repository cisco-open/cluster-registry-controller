// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

package controllers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func NewPredicateFuncs(filter func(meta metav1.Object, object runtime.Object, t string) bool) predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return filter(e.Meta, e.Object, "create")
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return filter(e.MetaNew, e.ObjectNew, "update")
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return filter(e.Meta, e.Object, "delete")
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return filter(e.Meta, e.Object, "generic")
		},
	}
}
