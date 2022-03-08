// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package clusters

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ManagedReconciler interface {
	reconcile.Reconciler

	PreCheck(ctx context.Context, client client.Client) error
	DoCleanup()
	GetName() string
	GetManager() ctrl.Manager
	SetManager(mgr ctrl.Manager)
	GetRecorder() record.EventRecorder
	GetContext() context.Context
	SetContext(ctx context.Context)
	GetLogger() logr.Logger
	SetLogger(l logr.Logger)
	GetClient() client.Client
	SetClient(client client.Client)
	Start(ctx context.Context) error
	SetScheme(scheme *runtime.Scheme)
	SetupWithController(ctx context.Context, ctrl controller.Controller) error
	SetupWithManager(ctx context.Context, mgr ctrl.Manager) error
}

type ManagedReconcilerBase struct {
	name string
	log  logr.Logger

	ctx      context.Context
	mgr      ctrl.Manager
	recorder record.EventRecorder
	scheme   *runtime.Scheme
	client   client.Client
}

func NewManagedReconciler(name string, log logr.Logger) ManagedReconciler {
	return &ManagedReconcilerBase{
		name: name,
		log:  log,
	}
}

func (r *ManagedReconcilerBase) GetClient() client.Client {
	return r.client
}

func (r *ManagedReconcilerBase) SetClient(client client.Client) {
	r.client = client
}

func (r *ManagedReconcilerBase) PreCheck(ctx context.Context, client client.Client) error {
	return nil
}

func (r *ManagedReconcilerBase) Start(ctx context.Context) error {
	return nil
}

func (r *ManagedReconcilerBase) DoCleanup() {
	r.ctx = nil
	r.mgr = nil
	r.recorder = nil
	r.scheme = nil
}

func (r *ManagedReconcilerBase) GetRecorder() record.EventRecorder {
	return r.recorder
}

func (r *ManagedReconcilerBase) GetName() string {
	return r.name
}

func (r *ManagedReconcilerBase) GetLogger() logr.Logger {
	return r.log
}

func (r *ManagedReconcilerBase) SetLogger(l logr.Logger) {
	r.log = l
}

func (r *ManagedReconcilerBase) GetContext() context.Context {
	return r.ctx
}

func (r *ManagedReconcilerBase) SetContext(ctx context.Context) {
	if r.ctx != nil {
		return
	}

	r.ctx = ctx
}

func (r *ManagedReconcilerBase) GetManager() ctrl.Manager {
	return r.mgr
}

func (r *ManagedReconcilerBase) SetManager(mgr ctrl.Manager) {
	if r.mgr != nil {
		return
	}

	r.mgr = mgr
	r.SetScheme(mgr.GetScheme())
	r.recorder = mgr.GetEventRecorderFor(r.name)
}

func (r *ManagedReconcilerBase) SetScheme(scheme *runtime.Scheme) {
	r.scheme = scheme
}

func (r *ManagedReconcilerBase) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *ManagedReconcilerBase) SetupWithController(ctx context.Context, ctrl controller.Controller) error {
	r.SetContext(ctx)

	return nil
}

func (r *ManagedReconcilerBase) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	r.SetManager(mgr)

	return nil
}
