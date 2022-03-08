// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"emperror.dev/errors"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/go-logr/logr"
	"github.com/throttled/throttled"
	authorizationv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/resources"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/clusters"
	"wwwin-github.cisco.com/cisco-app-networking/cluster-registry-controller/pkg/util"
)

const (
	originalNameLabel = "cluster-registry.k8s.cisco.com/original-name"
)

type syncReconciler struct {
	clusters.ManagedReconciler

	gvk             schema.GroupVersionKind
	localGVK        schema.GroupVersionKind
	localClusterID  string
	localMgr        ctrl.Manager
	localRecorder   record.EventRecorder
	clustersManager *clusters.Manager
	rateLimiter     throttled.RateLimiter

	clusterID      string
	ctrl           controller.Controller
	queue          workqueue.RateLimitingInterface
	rule           *clusterregistryv1alpha1.ResourceSyncRule
	localInformers map[string]struct{}

	localClient client.Client
	localCache  cache.Cache

	resourceNameMutated bool
}

type SyncReconcilerOption func(r *syncReconciler)

func WithRateLimiter(rateLimiter throttled.RateLimiter) SyncReconcilerOption {
	return func(r *syncReconciler) {
		r.rateLimiter = rateLimiter
	}
}

func NewSyncReconciler(name string, localMgr ctrl.Manager, rule *clusterregistryv1alpha1.ResourceSyncRule, log logr.Logger, clusterID string, clustersManager *clusters.Manager, opts ...SyncReconcilerOption) (SyncReconciler, error) {
	r := &syncReconciler{
		ManagedReconciler: clusters.NewManagedReconciler(name, log),

		gvk:             schema.GroupVersionKind(rule.Spec.GVK),
		localMgr:        localMgr,
		localRecorder:   localMgr.GetEventRecorderFor("cluster-controller"),
		clustersManager: clustersManager,
		rule:            rule,
		clusterID:       clusterID,
		localInformers:  make(map[string]struct{}),
	}

	_, r.localGVK = clusterregistryv1alpha1.MatchedRules(rule.Spec.Rules).GetMutatedGVK(r.gvk)

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

func (r *syncReconciler) PreCheck(ctx context.Context, client client.Client) error {
	for _, verb := range []string{"get", "list", "watch"} {
		attr := &authorizationv1.ResourceAttributes{
			Verb:     verb,
			Group:    r.gvk.Group,
			Version:  r.gvk.Version,
			Resource: strings.ToLower(pluralize.NewClient().Plural(r.gvk.Kind)),
		}
		selfSubjectAccessReview := authorizationv1.SelfSubjectAccessReview{
			Spec: authorizationv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: attr,
			},
		}

		err := client.Create(ctx, &selfSubjectAccessReview)
		if err != nil {
			return errors.WrapIfWithDetails(err, "failed to create self subject access review", "attributes", attr)
		}

		if !selfSubjectAccessReview.Status.Allowed {
			return errors.Errorf("do not have access to %s gvk: %s", attr.Verb, r.gvk)
		}
	}

	if err := r.writePreCheck(ctx); err != nil {
		return err
	}

	return nil
}

// WritePreCheck Check for write permissions on local cluster
func (r *syncReconciler) writePreCheck(ctx context.Context) error {
	localClient, err := client.New(r.localMgr.GetConfig(), client.Options{
		Scheme: r.localMgr.GetScheme(),
		Mapper: r.localMgr.GetRESTMapper(),
	})
	if err != nil {
		return errors.WrapIfWithDetails(err, "error creating local client")
	}

	for _, verb := range []string{"create", "patch", "update", "delete"} {
		attr := &authorizationv1.ResourceAttributes{
			Verb:     verb,
			Group:    r.localGVK.Group,
			Version:  r.localGVK.Version,
			Resource: strings.ToLower(pluralize.NewClient().Plural(r.localGVK.Kind)),
		}
		selfSubjectAccessReview := authorizationv1.SelfSubjectAccessReview{
			Spec: authorizationv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: attr,
			},
		}

		err = localClient.Create(ctx, &selfSubjectAccessReview)
		if err != nil {
			return errors.WrapIfWithDetails(err, "failed to create self subject access review", "attributes", attr)
		}

		if !selfSubjectAccessReview.Status.Allowed {
			return errors.Errorf("do not have local access to %s gvk: %s", attr.Verb, r.gvk)
		}
	}

	return nil
}

func (r *syncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	result, err := r.reconcile(ctx, req)
	if err != nil {
		r.localRecorder.Event(r.rule, corev1.EventTypeWarning, "ObjectNotReconciled", fmt.Sprintf("could not reconcile (resource: %s): %s", req, err.Error()))

		return result, err
	}

	return result, nil
}

func (r *syncReconciler) initObjectFromGVK(gvk schema.GroupVersionKind) client.Object {
	var object client.Object
	obj, err := r.localClient.Scheme().New(gvk)
	if err != nil {
		object = &unstructured.Unstructured{}
		object.GetObjectKind().SetGroupVersionKind(gvk)
	} else {
		object = obj.(client.Object) // nolint:forcetypeassert
	}

	return object
}

func (r *syncReconciler) reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.GetLogger().WithValues("resource", req.NamespacedName)

	obj := r.initObjectFromGVK(r.gvk)
	obj.SetName(req.Name)
	obj.SetNamespace(req.Namespace)

	// check namespace existence
	if req.Namespace != "" {
		err := r.localClient.Get(ctx, types.NamespacedName{
			Name: req.Namespace,
		}, &corev1.Namespace{})
		if err != nil && !apierrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}

		if apierrors.IsNotFound(err) {
			msg := "namespace does not exists locally"
			r.localRecorder.Event(r.rule, corev1.EventTypeWarning, "ObjectNotReconciledMissingNamespace", fmt.Sprintf("could not reconcile (resource: %s): %s", req, msg))
			log.Info(msg)

			return ctrl.Result{
				RequeueAfter: time.Second * 30, //nolint:gomnd
			}, nil
		}
	}

	err := r.GetClient().Get(ctx, req.NamespacedName, obj)
	if apierrors.IsNotFound(err) || err == nil && !obj.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, r.deleteResource(ctx, obj, log)
	}
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "could not get object")
	}

	if r.rateLimiter != nil {
		limited, _, err := r.rateLimiter.RateLimit(req.String(), 1)
		if err != nil {
			return ctrl.Result{}, errors.WrapIf(err, "could not rate limit")
		}
		if limited {
			msg := "ratelimited, too frequent reconciles were happening for this object"
			r.localRecorder.Event(r.rule, corev1.EventTypeWarning, "ObjectReconcileRateLimited", fmt.Sprintf("%s (resource: %s)", msg, req))
			log.Info(msg)

			return ctrl.Result{
				RequeueAfter: time.Second * 30, // nolint:gomnd
			}, nil
		}
	}

	ok, matchedRules, err := r.rule.Match(obj)
	if !ok {
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "could not match object")
	}

	log.Info("reconciling", "gvk", r.gvk)

	obj, err = r.mutateObject(obj, matchedRules)
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "could not mutate object")
	}

	rec := reconciler.NewGenericReconciler(
		r.localClient,
		log,
		reconciler.ReconcilerOpts{
			EnableRecreateWorkloadOnImmutableFieldChange: true,
			Scheme: r.localClient.Scheme(),
		},
	)

	var desiredObject client.Object
	if desiredObject, ok = obj.DeepCopyObject().(client.Object); !ok {
		return ctrl.Result{}, errors.New("invalid object")
	}

	_, err = rec.ReconcileResource(obj, r.getObjectDesiredState())
	if apierrors.IsAlreadyExists(errors.Cause(err)) {
		log.Info("object already exists, requeue")

		return ctrl.Result{
			Requeue: true,
		}, nil
	}
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "could not reconcile object")
	}
	log.Info("object reconciled")

	err = r.localClient.Get(ctx, client.ObjectKey{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}, obj)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "could not get object")
	}

	if matchedRules.GetMutationSyncStatus() {
		desiredObject.SetResourceVersion(obj.GetResourceVersion())
		err = r.localClient.Status().Update(ctx, desiredObject)
		if err != nil {
			return ctrl.Result{}, errors.WrapIf(err, "could not update object status")
		}
	}

	if r.rule.UID != "" {
		r.localRecorder.Event(r.rule, corev1.EventTypeNormal, "ObjectReconciled", fmt.Sprintf("object reconciled (resource: %s)", req))
	}

	return ctrl.Result{}, nil
}

func (r *syncReconciler) Start(ctx context.Context) error {
	// set local cluster id
	if r.localClusterID == "" {
		localClusterID, err := GetClusterID(ctx, r.localClient)
		if err != nil {
			return errors.WrapIf(err, "could not get local cluster id")
		}
		r.localClusterID = string(localClusterID)
		r.GetLogger().Info("set local cluster id", "id", r.localClusterID)
	}

	// init local informer
	_, gvk := clusterregistryv1alpha1.MatchedRules(r.rule.Spec.Rules).GetMutatedGVK(schema.GroupVersionKind(r.rule.Spec.GVK))
	obj := r.initObjectFromGVK(gvk)
	err := r.initLocalInformer(ctx, obj)
	if err != nil {
		return errors.WithStackIf(err)
	}

	return nil
}

func (r *syncReconciler) SetupWithController(ctx context.Context, ctrl controller.Controller) error {
	err := r.ManagedReconciler.SetupWithController(ctx, ctrl)
	if err != nil {
		return err
	}

	localCache, err := r.createAndStartCache()
	if err != nil {
		return err
	}
	r.localCache = localCache

	localClient, err := r.createClient(r.localMgr.GetConfig(), localCache)
	if err != nil {
		return err
	}
	r.localClient = localClient

	isObjectMatch := func(obj client.Object, gvk schema.GroupVersionKind) bool {
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		ok, _, err := r.rule.Match(obj)
		if err != nil {
			r.GetLogger().Error(err, "could not match object")

			return false
		}

		return ok
	}

	gvk := schema.GroupVersionKind(r.rule.Spec.GVK)
	obj := r.initObjectFromGVK(gvk)

	// set watcher for gvk
	err = ctrl.Watch(
		&source.Kind{
			Type: obj,
		},
		handler.EnqueueRequestsFromMapFunc(func(obj client.Object) []reconcile.Request {
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Name:      obj.GetName(),
						Namespace: obj.GetNamespace(),
					},
				},
			}
		}),
		predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				if r.isOwnedByUs(e.Object) {
					return false
				}

				return isObjectMatch(e.Object, gvk)
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				oldRV := e.ObjectOld.GetResourceVersion()
				e.ObjectOld.SetResourceVersion(e.ObjectNew.GetResourceVersion())
				defer e.ObjectOld.SetResourceVersion(oldRV)

				options := []patch.CalculateOption{
					reconciler.IgnoreManagedFields(),
				}

				patchResult, err := patch.DefaultPatchMaker.Calculate(e.ObjectOld, e.ObjectNew, options...)
				if err != nil {
					return true
				} else if patchResult.IsEmpty() {
					return false
				}

				if r.isOwnedByUs(e.ObjectNew) {
					return false
				}

				return isObjectMatch(e.ObjectNew, gvk)
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				if r.isOwnedByUs(e.Object) {
					return false
				}

				return isObjectMatch(e.Object, gvk)
			},
			GenericFunc: func(e event.GenericEvent) bool {
				if r.isOwnedByUs(e.Object) {
					return false
				}

				return isObjectMatch(e.Object, gvk)
			},
		},
	)
	if err != nil {
		return err
	}

	err = ctrl.Watch(&InMemorySource{
		reconciler: r,
	}, handler.Funcs{})
	if err != nil {
		return err
	}

	r.ctrl = ctrl

	return nil
}

func (r *syncReconciler) GetRule() *clusterregistryv1alpha1.ResourceSyncRule {
	return r.rule
}

func (r *syncReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	return errors.New("not implemented")
}

func (r *syncReconciler) mutateObject(current client.Object, matchedRules clusterregistryv1alpha1.MatchedRules) (client.Object, error) {
	var ok bool

	var obj client.Object
	if obj, ok = current.DeepCopyObject().(client.Object); !ok {
		return nil, errors.New("invalid object")
	}

	objAnnotations := obj.GetAnnotations()
	if objAnnotations == nil {
		objAnnotations = make(map[string]string)
	}

	for k, v := range matchedRules.GetMutationAnnotations().Add {
		objAnnotations[k] = v
	}

	for _, k := range matchedRules.GetMutationAnnotations().Remove {
		delete(objAnnotations, k)
	}

	objLabels := obj.GetLabels()
	if objLabels == nil {
		objLabels = make(map[string]string)
	}

	for k, v := range matchedRules.GetMutationLabels().Add {
		objLabels[k] = v
	}

	for _, k := range matchedRules.GetMutationLabels().Remove {
		delete(objLabels, k)
	}

	if objAnnotations[clusterregistryv1alpha1.OwnershipAnnotation] == "" {
		objAnnotations[clusterregistryv1alpha1.OwnershipAnnotation] = r.clusterID
	}

	if objLabels[clusterregistryv1alpha1.OwnershipAnnotation] == "" {
		objLabels[clusterregistryv1alpha1.OwnershipAnnotation] = r.clusterID
	}

	if mutated, gvk := matchedRules.GetMutatedGVK(obj.GetObjectKind().GroupVersionKind()); mutated {
		objAnnotations[clusterregistryv1alpha1.OriginalGVKAnnotation] = util.GVKToString(obj.GetObjectKind().GroupVersionKind())
		obj.GetObjectKind().SetGroupVersionKind(gvk)
	}

	delete(objAnnotations, patch.LastAppliedConfig)
	delete(objAnnotations, corev1.LastAppliedConfigAnnotation)
	obj.SetAnnotations(objAnnotations)
	obj.SetLabels(objLabels)

	obj.SetGeneration(0)
	obj.SetResourceVersion("")
	obj.SetUID("")
	obj.SetSelfLink("")
	obj.SetCreationTimestamp(metav1.Time{})
	obj.SetFinalizers(nil)
	obj.SetOwnerReferences(nil)
	obj.SetManagedFields(nil)

	if patches := matchedRules.GetMutationOverrides(); len(patches) > 0 { // nolint:nestif
		clusters, err := GetClusters(r.GetContext(), r.localClient)
		if err != nil {
			return nil, errors.WrapIf(err, "could not get clusters")
		}

		var localCluster clusterregistryv1alpha1.Cluster
		if localCluster, ok = clusters[types.UID(r.localClusterID)]; !ok {
			return nil, errors.NewWithDetails("could not find local cluster by id", "id", r.localClusterID)
		}

		syncedClusterID := types.UID(r.clusterID)
		if clusterID, ok := current.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]; ok {
			syncedClusterID = types.UID(clusterID)
		}
		var syncedCluster clusterregistryv1alpha1.Cluster
		if syncedCluster, ok = clusters[syncedClusterID]; !ok {
			return nil, errors.NewWithDetails("could not find synced cluster by id", "id", r.localClusterID)
		}

		modifiedPatches, err := util.K8SResourceOverlayPatchExecuteTemplates(patches, map[string]interface{}{
			"Object":       obj,
			"Cluster":      syncedCluster.DeepCopy(),
			"LocalCluster": localCluster.DeepCopy(),
		})
		if err != nil {
			return nil, errors.WrapIf(err, "could not execute templates on patches")
		}
		patches = modifiedPatches

		gvk := resources.ConvertGVK(obj.GetObjectKind().GroupVersionKind())
		patchFunc, err := resources.PatchYAMLModifier(resources.K8SResourceOverlay{
			GVK:     &gvk,
			Patches: patches,
		}, resources.NewObjectParser(r.localClient.Scheme()))
		if err != nil {
			return nil, errors.WrapIf(err, "could not get patch func for object")
		}

		patchedObject, err := patchFunc(obj)
		if err != nil {
			return nil, errors.WrapIf(err, "could not patch object")
		}

		if obj, ok = patchedObject.(client.Object); !ok {
			return nil, errors.New("invalid object")
		}
	}

	if current.GetName() != obj.GetName() {
		objLabels := obj.GetLabels()
		if objLabels == nil {
			objLabels = make(map[string]string)
		}
		objLabels[originalNameLabel] = current.GetName()
		obj.SetLabels(objLabels)
		r.resourceNameMutated = true
	}

	return obj, nil
}

func (r *syncReconciler) getObjectByOriginalName(ctx context.Context, obj client.Object, gvk schema.GroupVersionKind) (bool, client.Object, error) {
	var err error

	objects := &unstructured.UnstructuredList{}
	objects.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Kind:    fmt.Sprintf("%sList", gvk.Kind),
		Version: gvk.Version,
	})

	err = r.localClient.List(ctx, objects, client.InNamespace(obj.GetNamespace()), client.MatchingLabels(map[string]string{
		originalNameLabel: obj.GetName(),
		clusterregistryv1alpha1.OwnershipAnnotation: r.clusterID,
	}))
	if err != nil {
		return false, nil, err
	}

	if len(objects.Items) > 1 {
		return false, nil, errors.New("multiple renamed objects were found")
	}

	if len(objects.Items) == 1 {
		return true, objects.Items[0].DeepCopy(), nil
	}

	return false, nil, nil
}

func (r *syncReconciler) deleteResource(ctx context.Context, obj client.Object, log logr.Logger) error {
	var object, current client.Object
	var ok bool
	if object, ok = obj.DeepCopyObject().(client.Object); !ok {
		return errors.New("invalid object")
	}

	isGVKMutated := r.gvk != r.localGVK
	if isGVKMutated {
		object.GetObjectKind().SetGroupVersionKind(r.localGVK)
	}

	if current, ok = object.DeepCopyObject().(client.Object); !ok {
		return errors.New("invalid object")
	}

	err := r.localClient.Get(ctx, types.NamespacedName{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}, current)
	if apierrors.IsNotFound(err) { // nolint:nestif
		if !r.resourceNameMutated { // already deleted
			return nil
		}
		if ok, obj, err := r.getObjectByOriginalName(ctx, obj, current.GetObjectKind().GroupVersionKind()); err != nil {
			return err
		} else if ok {
			current = obj
		} else {
			return nil
		}
	} else if err != nil {
		return err
	}

	log = log.WithValues("resource", types.NamespacedName{
		Name:      current.GetName(),
		Namespace: current.GetNamespace(),
	})

	ownerClusterID := current.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]

	if ownerClusterID == "" {
		log.V(1).Info("deletion is skipped, object is owned by this cluster")

		return nil
	}

	if r.isOwnedByAnotherAliveCluster(ownerClusterID) {
		log.V(1).Info("deletion is skipped, owned by another live cluster")

		return nil
	}

	err = r.localClient.Delete(ctx, current)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	log.Info("object deleted")

	return nil
}

func (r *syncReconciler) isOwnedByAnotherAliveCluster(ownerClusterID string) bool {
	return ownerClusterID != "" && r.clustersManager.GetAliveClustersByID()[ownerClusterID] != nil && ownerClusterID != r.clusterID
}

func (r *syncReconciler) isOwnedByUs(object client.Object) bool {
	return object.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation] == r.localClusterID
}

func (r *syncReconciler) getObjectDesiredState() *reconciler.DynamicDesiredState {
	return &reconciler.DynamicDesiredState{
		BeforeUpdateFunc: func(current, desired runtime.Object) error {
			for _, f := range []func(current, desired runtime.Object) error{
				reconciler.ServiceIPModifier,
			} {
				err := f(current, desired)
				if err != nil {
					return err
				}
			}

			return nil
		},
		ShouldCreateFunc: func(desired runtime.Object) (bool, error) {
			metaObj, err := meta.Accessor(desired)
			if err != nil {
				return false, err
			}

			// sync disabled for this resource
			if _, ok := metaObj.GetAnnotations()[clusterregistryv1alpha1.SyncDisabledAnnotation]; ok {
				return false, nil
			}

			ownerClusterID := metaObj.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]
			// the resource is coming from through an intermediary but marked as owned by this cluster
			if ownerClusterID != "" && r.localClusterID == ownerClusterID {
				return false, nil
			}

			// this resource is owned by another live cluster - sync allowed only from that cluster
			if r.isOwnedByAnotherAliveCluster(ownerClusterID) {
				return false, nil
			}

			return true, nil
		},
		ShouldUpdateFunc: func(current, desired runtime.Object) (bool, error) {
			metaObj, err := meta.Accessor(current)
			if err != nil {
				return false, err
			}

			// sync disabled for this resource
			if _, ok := metaObj.GetAnnotations()[clusterregistryv1alpha1.SyncDisabledAnnotation]; ok {
				return false, nil
			}

			// this resources is owned by this cluster
			ownerClusterID := metaObj.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]
			if ownerClusterID == "" {
				return false, nil
			}

			// this resource is owned by another live cluster - sync is only allowed from that cluster
			if r.isOwnedByAnotherAliveCluster(ownerClusterID) {
				return false, nil
			}

			return true, nil
		},
	}
}

func (r *syncReconciler) setQueue(q workqueue.RateLimitingInterface) {
	r.queue = q
}

func (r *syncReconciler) initLocalInformer(ctx context.Context, obj client.Object) error {
	key := obj.GetObjectKind().GroupVersionKind().String()

	if _, ok := r.localInformers[key]; ok {
		return nil
	}

	r.GetLogger().Info("init local informer", "gvk", key)

	localInformer, err := r.localCache.GetInformer(ctx, obj)
	if err != nil {
		return errors.WrapIf(err, "could not create local informer for clusters")
	}

	err = r.ctrl.Watch(&source.Informer{
		Informer: localInformer,
	}, handler.EnqueueRequestsFromMapFunc(func(obj client.Object) []reconcile.Request {
		name := obj.GetName()
		if originalName, ok := obj.GetLabels()[originalNameLabel]; ok {
			name = originalName
		}

		return []reconcile.Request{
			{
				NamespacedName: types.NamespacedName{
					Name:      name,
					Namespace: obj.GetNamespace(),
				},
			},
		}
	}), r.localPredicate())
	if err != nil {
		return errors.WrapIf(err, "could not create watch for local informer")
	}

	r.localInformers[key] = struct{}{}

	return nil
}

func (r *syncReconciler) localPredicate() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			_, ok := e.Object.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]

			return ok
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			_, ok := e.ObjectOld.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]
			if !ok {
				return ok
			}

			oldRV := e.ObjectOld.GetResourceVersion()
			e.ObjectOld.SetResourceVersion(e.ObjectNew.GetResourceVersion())
			defer e.ObjectOld.SetResourceVersion(oldRV)

			options := []patch.CalculateOption{
				reconciler.IgnoreManagedFields(),
			}

			patchResult, err := patch.DefaultPatchMaker.Calculate(e.ObjectOld, e.ObjectNew, options...)
			if err != nil {
				return true
			} else if patchResult.IsEmpty() {
				return false
			}

			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			_, ok := e.Object.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]

			return ok
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
}

func (r *syncReconciler) createClient(config *rest.Config, cache cache.Cache) (client.Client, error) {
	cli, err := client.New(config, client.Options{
		Scheme: r.localMgr.GetScheme(),
		Mapper: r.localMgr.GetRESTMapper(),
	})
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

func (r *syncReconciler) createAndStartCache() (cache.Cache, error) {
	cche, err := cache.New(r.localMgr.GetConfig(), cache.Options{
		Scheme: r.localMgr.GetScheme(),
		Mapper: r.localMgr.GetRESTMapper(),
	})
	if err != nil {
		return nil, err
	}

	go func() {
		err = cche.Start(r.GetContext())
		if err != nil {
			r.GetLogger().Error(err, "could not start cache")
		}
		r.GetLogger().Info("cache stopped")
	}()

	if !cche.WaitForCacheSync(r.GetContext()) {
		return nil, errors.New("could not sync cache")
	}

	return cche, nil
}
