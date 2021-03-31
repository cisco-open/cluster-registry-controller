// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

package clusters

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

type ManagedController interface {
	SetLogger(l logr.Logger)
	GetReconciler() ManagedReconciler
	GetName() string
	Stop()
	Stopped() <-chan struct{}
	Start(ctx context.Context, mgr ctrl.Manager) error
}

type ManagedControllers map[string]ManagedController

type managedController struct {
	name              string
	reconciler        ManagedReconciler
	log               logr.Logger
	mgr               ctrl.Manager
	ctrl              controller.Controller
	ctrlContext       context.Context
	ctrlContextCancel context.CancelFunc
}

func NewManagedController(name string, r ManagedReconciler, l logr.Logger) ManagedController {
	return &managedController{
		name:       name,
		reconciler: r,
		log:        l.WithName(name),
	}
}

func (c *managedController) GetName() string {
	return c.name
}

func (c *managedController) GetReconciler() ManagedReconciler {
	return c.reconciler
}

func (c *managedController) SetLogger(l logr.Logger) {
	c.log = l
}

func (c *managedController) Stopped() <-chan struct{} {
	return c.ctrlContext.Done()
}

func (c *managedController) Stop() {
	if c.ctrlContextCancel != nil {
		c.ctrlContextCancel()
	}
}

func (c *managedController) Start(ctx context.Context, mgr ctrl.Manager) error {
	if c.ctrlContextCancel != nil {
		return nil
	}

	c.ctrlContext, c.ctrlContextCancel = context.WithCancel(ctx)
	c.mgr = mgr

	c.reconciler.SetManager(c.mgr)
	c.reconciler.SetLogger(c.log)

	var err error
	c.ctrl, err = controller.NewUnmanaged(c.GetName(), c.mgr, controller.Options{
		Reconciler: c.reconciler,
		Log:        c.log,
	})
	if err != nil {
		return err
	}

	err = c.reconciler.SetupWithController(c.ctrlContext, c.ctrl)
	if err != nil {
		return errors.WithStackIf(err)
	}

	// Start our controller in a goroutine so that we do not block.
	go func() {
		// Block until our controller manager is elected leader. We presume our entire
		// process will terminate if we lose leadership, so we don't need to handle that.
		<-mgr.Elected()

		started := mgr.GetCache().WaitForCacheSync(c.ctrlContext.Done())
		if !started {
			c.log.Error(err, "timeout while waiting for cache sync")

			return
		}

		// Start our controller. This will block until the stop channel is closed, or the
		// controller returns an error.
		c.log.Info("starting ctrl")
		if err := c.ctrl.Start(c.ctrlContext.Done()); err != nil {
			c.log.Error(err, "cannot run sync controller")
		}
		c.log.Info("ctrl stopped")
		<-c.ctrlContext.Done()
		c.reconciler.DoCleanup()
		c.ctrlContextCancel = nil
		c.ctrlContext = nil
		c.ctrl = nil
	}()

	return nil
}
