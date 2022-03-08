// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package clusters

import (
	"context"

	"emperror.dev/errors"
	"github.com/cenkalti/backoff"
	"github.com/go-logr/logr"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

type ManagedController interface {
	SetLogger(l logr.Logger)
	GetReconciler() ManagedReconciler
	GetName() string
	Stop()
	Stopped() <-chan struct{}
	Start(ctx context.Context, mgr ctrl.Manager) error
	GetRequiredClusterFeatures() []ClusterFeatureRequirement
	GetClient() client.Client
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
	client            client.Client
	cache             cache.Cache

	requiredClusterFeatures []ClusterFeatureRequirement
}

type ManagedControllerOption func(*managedController)

func WithRequiredClusterFeatures(features ...ClusterFeatureRequirement) ManagedControllerOption {
	return func(r *managedController) {
		if r.requiredClusterFeatures == nil {
			r.requiredClusterFeatures = make([]ClusterFeatureRequirement, 0)
		}
		r.requiredClusterFeatures = append(r.requiredClusterFeatures, features...)
	}
}

func NewManagedController(name string, r ManagedReconciler, l logr.Logger, options ...ManagedControllerOption) ManagedController {
	m := &managedController{
		name:                    name,
		reconciler:              r,
		log:                     l.WithName(name),
		requiredClusterFeatures: make([]ClusterFeatureRequirement, 0),
	}

	for _, o := range options {
		o(m)
	}

	return m
}

func (c *managedController) GetRequiredClusterFeatures() []ClusterFeatureRequirement {
	return c.requiredClusterFeatures
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
	c.ctrlContextCancel = nil
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

	check := func() error {
		err = c.reconciler.PreCheck(c.ctrlContext, c.mgr.GetClient())
		if err != nil {
			return errors.WrapIf(err, "pre check error")
		}

		err = c.start()
		if err != nil {
			return errors.WrapIf(err, "could not start controller")
		}

		return nil
	}

	go func() {
		err = check()
		if err != nil {
			c.log.Error(err, "")
		} else {
			return
		}
		ticker := backoff.NewTicker(backoff.NewExponentialBackOff())
		for {
			select {
			case <-ticker.C:
				err = check()
				if err != nil {
					c.log.Error(err, "")
				} else {
					return
				}
			case <-c.ctrlContext.Done():
				return
			}
		}
	}()

	return nil
}

func (c *managedController) GetClient() client.Client {
	return c.client
}

func (c *managedController) createClient(config *rest.Config, options client.Options, cache cache.Cache) (client.Client, error) {
	cli, err := client.New(config, options)
	if err != nil {
		return nil, err
	}

	cli, err = client.NewDelegatingClient(client.NewDelegatingClientInput{
		CacheReader:     cache,
		Client:          cli,
		UncachedObjects: nil,
	})
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func (c *managedController) createAndStartCache() (cache.Cache, error) {
	if c.ctrlContext == nil {
		return nil, errors.New("context is nil")
	}

	cche, err := cache.New(c.mgr.GetConfig(), cache.Options{
		Scheme: c.mgr.GetScheme(),
		Mapper: c.mgr.GetRESTMapper(),
	})
	if err != nil {
		return nil, err
	}

	go func() {
		err = cche.Start(c.ctrlContext)
		if err != nil {
			c.log.Error(err, "could not start cache")
		}
		c.log.Info("cache stopped")
	}()

	if !cche.WaitForCacheSync(c.ctrlContext) {
		return nil, errors.New("could not sync cache")
	}

	return cche, nil
}

func (c *managedController) getInjectFunc() inject.Func {
	return func(i interface{}) error {
		if _, ok := i.(inject.Cache); ok {
			if ok, err := inject.CacheInto(c.cache, i); err != nil {
				return err
			} else if ok {
				c.log.V(2).Info("cache injected", "i", i)

				return nil
			}

			return nil
		}

		if _, ok := i.(inject.Client); ok {
			if ok, err := inject.ClientInto(c.client, i); err != nil {
				return err
			} else if ok {
				c.log.V(2).Info("client injected", "i", i)

				return nil
			}

			return nil
		}

		return c.mgr.SetFields(i)
	}
}

func (c *managedController) start() error {
	var err error

	c.ctrl, err = controller.NewUnmanaged(c.GetName(), c.mgr, controller.Options{
		Reconciler: c.reconciler,
		Log:        c.log,
	})
	if err != nil {
		return errors.WithStackIf(err)
	}

	cache, err := c.createAndStartCache()
	if err != nil {
		return errors.WithStackIf(err)
	}
	c.cache = cache

	client, err := c.createClient(c.mgr.GetConfig(), client.Options{
		Scheme: c.mgr.GetScheme(),
		Mapper: c.mgr.GetRESTMapper(),
	}, cache)
	if err != nil {
		return errors.WithStackIf(err)
	}
	c.client = client

	if controller, ok := c.ctrl.(interface {
		InjectFunc(f inject.Func) error
	}); ok {
		err = controller.InjectFunc(c.getInjectFunc())
		if err != nil {
			return err
		}
	}

	c.reconciler.SetClient(c.client)

	err = c.reconciler.SetupWithController(c.ctrlContext, c.ctrl)
	if err != nil {
		return errors.WithStackIf(err)
	}

	// Start our controller in a goroutine so that we do not block.
	go func() {
		// Block until our controller manager is elected leader. We presume our entire
		// process will terminate if we lose leadership, so we don't need to handle that.
		<-c.mgr.Elected()

		started := c.mgr.GetCache().WaitForCacheSync(c.ctrlContext)
		if !started {
			c.log.Error(err, "timeout while waiting for cache sync")

			return
		}

		// Start our controller. This will block until the stop channel is closed, or the
		// controller returns an error.
		c.log.Info("starting ctrl")

		err = c.reconciler.Start(c.ctrlContext)
		if err != nil {
			c.log.Error(err, "")

			return
		}

		if err := c.ctrl.Start(c.ctrlContext); err != nil {
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
