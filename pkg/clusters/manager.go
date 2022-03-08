// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package clusters

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrClusterNotFound    = errors.New("cluster not found")
	ErrControllerNotFound = errors.New("controller not found")
)

type ManagerOption func(m *Manager)

type Manager struct {
	clusters map[string]*Cluster
	mu       *sync.RWMutex
	ctx      context.Context

	onBeforeAddFuncs    map[string]func(c *Cluster)
	onBeforeDeleteFuncs map[string]func(c *Cluster)
	onAfterAddFuncs     map[string]func(c *Cluster)
	onAfterDeleteFuncs  map[string]func()
}

func WithOnBeforeAddFunc(f func(c *Cluster), ids ...string) ManagerOption {
	return func(m *Manager) {
		m.AddOnBeforeAddFunc(f, ids...)
	}
}

func WithOnAfterAddFunc(f func(c *Cluster), ids ...string) ManagerOption {
	return func(m *Manager) {
		m.AddOnAfterAddFunc(f, ids...)
	}
}

func WithOnBeforeDeleteFunc(f func(c *Cluster), ids ...string) ManagerOption {
	return func(m *Manager) {
		m.AddOnBeforeDeleteFunc(f, ids...)
	}
}

func WithOnAfterDeleteFunc(f func(), ids ...string) ManagerOption {
	return func(m *Manager) {
		m.AddOnAfterDeleteFunc(f, ids...)
	}
}

func NewManager(ctx context.Context, options ...ManagerOption) *Manager {
	mgr := &Manager{
		clusters: make(map[string]*Cluster),
		mu:       &sync.RWMutex{},
		ctx:      ctx,
	}

	for _, opt := range options {
		opt(mgr)
	}

	go mgr.waitForStop(ctx)

	return mgr
}

func (m *Manager) AddOnBeforeAddFunc(f func(c *Cluster), ids ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.onBeforeAddFuncs == nil {
		m.onBeforeAddFuncs = make(map[string]func(c *Cluster))
	}

	id := time.Now().String()
	if len(ids) > 0 {
		id = ids[0]
	}
	m.onBeforeAddFuncs[id] = f
}

func (m *Manager) DeleteOnBeforeAddFunc(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.onBeforeAddFuncs, id)
}

func (m *Manager) AddOnAfterAddFunc(f func(c *Cluster), ids ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.onAfterAddFuncs == nil {
		m.onAfterAddFuncs = make(map[string]func(c *Cluster))
	}
	id := time.Now().String()
	if len(ids) > 0 {
		id = ids[0]
	}
	m.onAfterAddFuncs[id] = f
}

func (m *Manager) DeleteOnAfterAddFunc(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.onAfterAddFuncs, id)
}

func (m *Manager) AddOnBeforeDeleteFunc(f func(c *Cluster), ids ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.onBeforeDeleteFuncs == nil {
		m.onBeforeDeleteFuncs = make(map[string]func(c *Cluster))
	}
	id := time.Now().String()
	if len(ids) > 0 {
		id = ids[0]
	}
	m.onBeforeDeleteFuncs[id] = f
}

func (m *Manager) DeleteOnBeforeDeleteFunc(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.onBeforeDeleteFuncs, id)
}

func (m *Manager) AddOnAfterDeleteFunc(f func(), ids ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.onAfterDeleteFuncs == nil {
		m.onAfterDeleteFuncs = make(map[string]func())
	}
	id := time.Now().String()
	if len(ids) > 0 {
		id = ids[0]
	}
	m.onAfterDeleteFuncs[id] = f
}

func (m *Manager) DeleteOnAfterDeleteFunc(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.onAfterDeleteFuncs, id)
}

func (m *Manager) Exists(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.clusters[name]

	return ok
}

func (m *Manager) Get(name string) (*Cluster, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if cluster, ok := m.clusters[name]; ok {
		return cluster, nil
	}

	return nil, ErrClusterNotFound
}

func (m *Manager) GetAll() map[string]*Cluster {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.clusters
}

func (m *Manager) GetAliveClustersByID() map[string]*Cluster {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clusters := make(map[string]*Cluster)
	for _, cluster := range m.clusters {
		if cluster.IsAlive() {
			clusters[cluster.GetClusterID()] = cluster
		}
	}

	return clusters
}

func (m *Manager) Add(cluster *Cluster) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, f := range m.onBeforeAddFuncs {
		f(cluster)
	}

	m.clusters[cluster.GetName()] = cluster

	for _, f := range m.onAfterAddFuncs {
		f(cluster)
	}

	return nil
}

func (m *Manager) Remove(cluster *Cluster) error {
	if cluster == nil {
		return nil
	}

	if m.clusters[cluster.GetName()] == nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, f := range m.onBeforeDeleteFuncs {
		f(cluster)
	}

	cluster.Stop()

	delete(m.clusters, cluster.GetName())

	for _, f := range m.onAfterDeleteFuncs {
		f()
	}

	return nil
}

func (m *Manager) Stopped() <-chan struct{} {
	return m.ctx.Done()
}

func (m *Manager) waitForStop(ctx context.Context) {
	<-ctx.Done()

	for _, c := range m.clusters {
		c.Stop()
	}
}
