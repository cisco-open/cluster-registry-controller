// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package clusters

import (
	"context"
	"sync"
	"time"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/banzaicloud/operator-tools/pkg/resources"
)

const defaultLivenessCheckInterval = time.Second * 5

type Cluster struct {
	name      string
	k8sConfig *rest.Config

	ctrlOptions ctrl.Options

	ctx          context.Context
	ctxCancel    context.CancelFunc
	mgrCtx       context.Context
	mgrCtxCancel context.CancelFunc

	log logr.Logger
	mgr ctrl.Manager

	alive                 bool
	started               bool
	mgrStopped            bool
	secretID              *string
	clusterID             string
	livenessCheckInterval time.Duration
	onAliveFuncs          []ClusterFunc
	onDeadFuncs           []ClusterFunc
	features              map[string]ClusterFeature
	kubeconfig            []byte

	controllers        ManagedControllers
	pendingControllers ManagedControllers
	mu                 *sync.RWMutex
}

type (
	ClusterFunc func(c *Cluster) error
	Option      func(*Cluster)
)

func WithSecretID(secretID string) Option {
	return func(c *Cluster) {
		c.secretID = &secretID
	}
}

func WithKubeconfig(kubeconfig []byte) Option {
	return func(c *Cluster) {
		c.kubeconfig = kubeconfig
	}
}

func WithLivenessCheckInterval(interval time.Duration) Option {
	return func(c *Cluster) {
		c.livenessCheckInterval = interval
	}
}

func WithOnAliveFunc(f func(c *Cluster) error) Option {
	return func(c *Cluster) {
		c.AddOnAliveFunc(f)
	}
}

func WithOnDeadFunc(f func(c *Cluster) error) Option {
	return func(c *Cluster) {
		c.AddOnDeadFunc(f)
	}
}

func WithCtrlOption(options ctrl.Options) Option {
	return func(c *Cluster) {
		c.ctrlOptions = options
	}
}

func WithScheme(scheme *runtime.Scheme) Option {
	return func(c *Cluster) {
		c.ctrlOptions.Scheme = scheme
	}
}

func NewCluster(ctx context.Context, name string, k8sConfig *rest.Config, log logr.Logger, opts ...Option) (*Cluster, error) {
	c := &Cluster{
		name:      name,
		k8sConfig: k8sConfig,
		log:       log.WithName(name),
		ctrlOptions: ctrl.Options{
			Scheme:             runtime.NewScheme(),
			MetricsBindAddress: "0",
			Port:               0,
		},
		livenessCheckInterval: defaultLivenessCheckInterval,
		onAliveFuncs:          make([]ClusterFunc, 0),
		onDeadFuncs:           make([]ClusterFunc, 0),
		features:              make(map[string]ClusterFeature),

		controllers:        make(ManagedControllers),
		pendingControllers: make(ManagedControllers),
		mu:                 &sync.RWMutex{},
	}

	c.ctx, c.ctxCancel = context.WithCancel(ctx)

	// Loop through each option
	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Cluster) Start() error {
	if err := c.livenessCheck(); err != nil {
		return err
	}

	c.started = true

	go func(ctx context.Context, cluster *Cluster, log logr.Logger) {
		ticker := time.NewTicker(cluster.livenessCheckInterval)
		for {
			select {
			case <-ctx.Done():
				c.setDead()
				ticker.Stop()
				c.started = false
				log.V(2).Info("context cancelled")

				return
			case <-ticker.C:
				log.V(2).Info("liveness check")
				err := cluster.livenessCheck()
				if err != nil {
					log.V(2).Error(err, "")
				}
			}
		}
	}(c.ctx, c, c.log.WithName("liveness"))

	return nil
}

func (c *Cluster) AddFeature(feature ClusterFeature) {
	c.mu.Lock()
	c.features[feature.GetUID()] = feature
	c.mu.Unlock()

	// check pending controllers
	for _, controller := range c.pendingControllers {
		err := c.AddController(controller)
		if err != nil {
			c.log.Error(err, "could not add controller")
		}
	}
}

func (c *Cluster) RemoveFeature(uid string) {
	c.mu.Lock()
	delete(c.features, uid)
	c.mu.Unlock()

	// check controllers if they are still meet the feature requirements
	for _, controller := range c.controllers {
		if !c.checkRequiredClusterFeatures(controller) {
			c.RemoveController(controller)
			c.addPendingController(controller)
		}
	}
}

func (c *Cluster) AddOnAliveFunc(f ClusterFunc) {
	c.onAliveFuncs = append(c.onAliveFuncs, f)
	if c.IsAlive() {
		c.runClusterFunc(f)
	}
}

func (c *Cluster) AddOnDeadFunc(f ClusterFunc) {
	c.onDeadFuncs = append(c.onDeadFuncs, f)
}

func (c *Cluster) GetClusterID() string {
	return c.clusterID
}

func (c *Cluster) IsAlive() bool {
	return c.alive
}

func (c *Cluster) IsManagerRunning() bool {
	return !c.mgrStopped && c.mgr != nil
}

func (c *Cluster) Stop() {
	if c == nil {
		return
	}
	c.log.V(2).Info("shutdown cluster")
	c.ctxCancel()
}

func (c *Cluster) StopManager() {
	if c == nil {
		return
	}

	c.stopManager()
	c.mgrStopped = true
}

func (c *Cluster) stopManager() {
	if c == nil {
		return
	}

	if c.mgrCtx != nil && c.mgrCtxCancel != nil {
		c.log.V(2).Info("shutdown manager")
		c.mgrCtxCancel()
	}
}

func (c *Cluster) Stopped() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Cluster) ManagerStopped() <-chan struct{} {
	return c.mgrCtx.Done()
}

func (c *Cluster) GetName() string {
	return c.name
}

func (c *Cluster) GetManager() ctrl.Manager {
	return c.mgr
}

func (c *Cluster) GetSecretID() *string {
	return c.secretID
}

func (c *Cluster) GetKubeconfig() []byte {
	return c.kubeconfig
}

func (c *Cluster) AddController(controller ManagedController) error {
	if c.checkRequiredClusterFeatures(controller) {
		return c.addController(controller)
	}

	c.addPendingController(controller)

	return nil
}

func (c *Cluster) GetPendingControllers() ManagedControllers {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.pendingControllers
}

func (c *Cluster) addController(controller ManagedController) error {
	name := controller.GetName()

	if c.GetController(name) != nil {
		return nil
	}

	c.removePendingControllerByName(controller.GetName())

	c.mu.Lock()
	defer c.mu.Unlock()

	controller.SetLogger(c.log.WithName(name))
	c.controllers[name] = controller

	if c.IsManagerRunning() {
		err := controller.Start(c.mgrCtx, c.mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) HasController(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.controllers[name]

	return ok
}

func (c *Cluster) GetController(name string) ManagedController {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.controllers[name]
}

func (c *Cluster) GetControllers() ManagedControllers {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.controllers
}

func (c *Cluster) GetControllerByGVK(gvk resources.GroupVersionKind) ManagedController {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for name, controller := range c.controllers {
		if name == gvk.Group {
			return controller
		}
	}

	return nil
}

func (c *Cluster) RemoveController(controller ManagedController) {
	c.RemoveControllerByName(controller.GetName())
}

func (c *Cluster) RemoveControllerByName(name string) {
	c.removePendingControllerByName(name)

	controller := c.GetController(name)

	c.mu.Lock()
	defer c.mu.Unlock()

	if controller == nil {
		return
	}

	controller.Stop()
	delete(c.controllers, name)
}

func (c *Cluster) setAlive() {
	if c.alive {
		return
	}

	if !c.mgrStopped {
		err := c.StartManager()
		if err != nil {
			c.log.Error(err, "")
		}
	}

	for _, f := range c.onAliveFuncs {
		c.runClusterFunc(f)
	}

	c.alive = true
}

func (c *Cluster) setDead() {
	if !c.alive {
		return
	}

	c.stopManager()

	for _, f := range c.onDeadFuncs {
		c.runClusterFunc(f)
	}

	c.alive = false
}

func (c *Cluster) livenessCheck() error {
	clientset, err := kubernetes.NewForConfig(c.k8sConfig)
	if err != nil {
		c.setDead()

		return errors.WrapIf(err, "could not get cluster ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ns, err := clientset.CoreV1().Namespaces().Get(ctx, metav1.NamespaceSystem, metav1.GetOptions{})
	if err != nil {
		c.setDead()

		return errors.WrapIf(err, "could not get cluster ID")
	}

	c.setAlive()
	c.clusterID = string(ns.UID)

	return nil
}

func (c *Cluster) StartManager() error {
	var err error

	if c.mgrCtx != nil {
		return nil
	}

	c.mgrStopped = false
	c.mgrCtx, c.mgrCtxCancel = context.WithCancel(c.ctx)

	c.mgr, err = ctrl.NewManager(c.k8sConfig, c.ctrlOptions)
	if err != nil {
		return errors.WrapIf(err, "could not create manager")
	}

	go func() {
		err = c.mgr.Start(c.mgrCtx)
		if err != nil {
			c.log.Error(err, "could not start manager")
		} else {
			c.log.V(2).Info("manager stopped")
		}
		<-c.mgrCtx.Done()
		c.mgrCtx = nil
		c.mgrCtxCancel = nil
		c.mgr = nil
	}()

	c.mgr.GetCache().WaitForCacheSync(c.mgrCtx)

	c.log.V(2).Info("manager started")

	if c.controllers != nil {
		for _, mctrl := range c.controllers {
			c.log.V(2).Info("start ctrl", "name", mctrl.GetName())
			err := mctrl.Start(c.mgrCtx, c.mgr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Cluster) runClusterFunc(f ClusterFunc) {
	//nolint:gomnd
	backoff := wait.Backoff{
		Steps:    5,
		Duration: 100 * time.Millisecond,
		Factor:   1.5,
		Jitter:   0.1,
	}
	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		if err := f(c); err != nil {
			c.log.Error(err, "could not run onAlive function")

			return false, nil
		}

		return true, nil
	})
	if err != nil {
		c.log.Error(err, "could not run onAlive function")
	}
}

func (c *Cluster) addPendingController(controller ManagedController) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.pendingControllers[controller.GetName()] = controller
}

func (c *Cluster) removePendingControllerByName(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.pendingControllers, name)
}

func (c *Cluster) checkRequiredClusterFeatures(controller ManagedController) bool {
	var requiredFeatures []ClusterFeatureRequirement

	if requiredFeatures = controller.GetRequiredClusterFeatures(); len(requiredFeatures) == 0 {
		return true
	}

	for _, rf := range requiredFeatures {
		ok := rf.Match(c.features)
		if !ok {
			return false
		}
	}

	return true
}
