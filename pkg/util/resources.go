// Copyright (c) 2019 Banzai Cloud Zrt. All Rights Reserved.

package util

import (
	"k8s.io/apimachinery/pkg/runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
)

type DynamicDesiredState struct {
	DesiredState     reconciler.DesiredState
	BeforeCreateFunc func(desired runtime.Object) error
	BeforeUpdateFunc func(current, desired runtime.Object) error
	BeforeDeleteFunc func(current runtime.Object) error
	CreateOptions    []runtimeClient.CreateOption
	UpdateOptions    []runtimeClient.UpdateOption
	DeleteOptions    []runtimeClient.DeleteOption
	ShouldCreateFunc func(desired runtime.Object) (bool, error)
	ShouldUpdateFunc func(current, desired runtime.Object) (bool, error)
	ShouldDeleteFunc func(desired runtime.Object) (bool, error)
}

func (s DynamicDesiredState) GetDesiredState() reconciler.DesiredState {
	return s.DesiredState
}

func (s DynamicDesiredState) ShouldCreate(desired runtime.Object) (bool, error) {
	if s.ShouldCreateFunc != nil {
		return s.ShouldCreateFunc(desired)
	}

	return true, nil
}

func (s DynamicDesiredState) ShouldUpdate(current, desired runtime.Object) (bool, error) {
	if s.ShouldUpdateFunc != nil {
		return s.ShouldUpdateFunc(current, desired)
	}

	return true, nil
}

func (s DynamicDesiredState) ShouldDelete(desired runtime.Object) (bool, error) {
	if s.ShouldDeleteFunc != nil {
		return s.ShouldDeleteFunc(desired)
	}

	return true, nil
}

func (s DynamicDesiredState) BeforeCreate(desired runtime.Object) error {
	if s.BeforeCreateFunc != nil {
		return s.BeforeCreateFunc(desired)
	}

	return nil
}

func (s DynamicDesiredState) BeforeUpdate(current, desired runtime.Object) error {
	if s.BeforeUpdateFunc != nil {
		return s.BeforeUpdateFunc(current, desired)
	}

	return nil
}

func (s DynamicDesiredState) BeforeDelete(current runtime.Object) error {
	if s.BeforeDeleteFunc != nil {
		return s.BeforeDeleteFunc(current)
	}

	return nil
}

func (s DynamicDesiredState) GetCreateptions() []runtimeClient.CreateOption {
	return s.CreateOptions
}

func (s DynamicDesiredState) GetUpdateOptions() []runtimeClient.UpdateOption {
	return s.UpdateOptions
}

func (s DynamicDesiredState) GetDeleteOptions() []runtimeClient.DeleteOption {
	return s.DeleteOptions
}
