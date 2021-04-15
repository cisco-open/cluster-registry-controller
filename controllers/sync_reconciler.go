// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

package controllers

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	"github.com/throttled/throttled"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/banzaicloud/cluster-registry-controller/pkg/clusters"
	"github.com/banzaicloud/cluster-registry-controller/pkg/util"
	clusterregistryv1alpha1 "github.com/banzaicloud/cluster-registry/api/v1alpha1"
	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/resources"
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

func (r *syncReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	result, err := r.reconcile(req)
	if err != nil {
		r.localRecorder.Event(r.rule, corev1.EventTypeWarning, "ObjectNotReconciled", fmt.Sprintf("could not reconcile (resource: %s): %s", req, err.Error()))

		return result, err
	}

	if r.rule.UID != "" {
		r.localRecorder.Event(r.rule, corev1.EventTypeNormal, "ObjectReconciled", fmt.Sprintf("object reconciled (resource: %s)", req))
	}

	return result, nil
}

func (r *syncReconciler) reconcile(req ctrl.Request) (ctrl.Result, error) {
	log := r.GetLogger().WithValues("resource", req.NamespacedName)

	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(r.gvk)
	obj.SetName(req.Name)
	obj.SetNamespace(req.Namespace)

	log.Info("reconciling", "gvk", r.gvk)

	// check namespace existence
	if req.Namespace != "" {
		err := r.localMgr.GetClient().Get(r.GetContext(), types.NamespacedName{
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

	err := r.GetManager().GetClient().Get(r.GetContext(), req.NamespacedName, obj)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, r.deleteResource(obj, log)
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

			return ctrl.Result{}, nil
		}
	}

	ok, matchedRules, err := r.rule.Match(obj)
	if !ok {
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "could not match object")
	}

	obj, err = r.mutateObject(obj, matchedRules)
	if err != nil {
		return ctrl.Result{}, errors.WrapIf(err, "could not mutate object")
	}

	rec := reconciler.NewGenericReconciler(
		r.localMgr.GetClient(),
		log,
		reconciler.ReconcilerOpts{
			EnableRecreateWorkloadOnImmutableFieldChange: true,
			Scheme: r.localMgr.GetScheme(),
		},
	)

	desiredObject := obj.DeepCopy()
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

	err = r.localMgr.GetClient().Get(r.GetContext(), client.ObjectKey{
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
		err = r.localMgr.GetClient().Status().Update(r.GetContext(), desiredObject)
		if err != nil {
			return ctrl.Result{}, errors.WrapIf(err, "could not update object status")
		}
	}

	return ctrl.Result{}, nil
}

func (r *syncReconciler) Start() error {
	// set local cluster id
	if r.localClusterID == "" {
		localClusterID, err := GetClusterID(r.GetContext(), r.localMgr.GetClient())
		if err != nil {
			return errors.WrapIf(err, "could not get local cluster id")
		}
		r.localClusterID = string(localClusterID)
		r.GetLogger().Info("set local cluster id", "id", r.localClusterID)
	}

	// init local informer
	_, gvk := clusterregistryv1alpha1.MatchedRules(r.rule.Spec.Rules).GetMutatedGVK(schema.GroupVersionKind(r.rule.Spec.GVK))
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)

	err := r.initLocalInformer(r.GetContext(), obj)
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

	gvk := schema.GroupVersionKind(r.rule.Spec.GVK)

	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)

	// set watcher for gvk
	err = ctrl.Watch(&source.Kind{
		Type: obj,
	}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(obj handler.MapObject) []reconcile.Request {
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Name:      obj.Meta.GetName(),
						Namespace: obj.Meta.GetNamespace(),
					},
				},
			}
		}),
	}, predicate.NewPredicateFuncs(func(meta metav1.Object, object runtime.Object) bool {
		object.GetObjectKind().SetGroupVersionKind(gvk)
		ok, _, err := r.rule.Match(object)
		if err != nil {
			r.GetLogger().Error(err, "could not match object")

			return false
		}

		return ok
	}))
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
	return r.ManagedReconciler.SetupWithManager(ctx, mgr)
}

func (r *syncReconciler) mutateObject(current *unstructured.Unstructured, matchedRules clusterregistryv1alpha1.MatchedRules) (*unstructured.Unstructured, error) {
	var ok bool

	obj := current.DeepCopy()

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

	if mutated, gvk := matchedRules.GetMutatedGVK(obj.GetObjectKind().GroupVersionKind()); mutated {
		objAnnotations[clusterregistryv1alpha1.OriginalGVKAnnotation] = util.GVKToString(obj.GetObjectKind().GroupVersionKind())
		obj.GetObjectKind().SetGroupVersionKind(gvk)
	}

	delete(objAnnotations, patch.LastAppliedConfig)
	delete(objAnnotations, corev1.LastAppliedConfigAnnotation)
	obj.SetAnnotations(objAnnotations)

	obj.SetGeneration(0)
	obj.SetResourceVersion("")
	obj.SetUID("")
	obj.SetSelfLink("")
	obj.SetCreationTimestamp(metav1.Time{})
	obj.SetFinalizers(nil)

	gvk := resources.ConvertGVK(obj.GetObjectKind().GroupVersionKind())
	patchFunc, err := resources.PatchYAMLModifier(resources.K8SResourceOverlay{
		GVK:     &gvk,
		Patches: matchedRules.GetMutationOverrides(),
	}, resources.NewObjectParser(runtime.NewScheme()))
	if err != nil {
		return nil, errors.WrapIf(err, "could not get patch func for object")
	}

	patchedObject, err := patchFunc(obj)
	if err != nil {
		return nil, errors.WrapIf(err, "could not patch object")
	}

	if obj, ok = patchedObject.(*unstructured.Unstructured); !ok {
		return nil, errors.New("invalid object")
	}

	return obj, nil
}

func (r *syncReconciler) deleteResource(obj *unstructured.Unstructured, log logr.Logger) error {
	log.Info("object was removed, trying to delete locally as well")

	object := obj.DeepCopy()
	object.SetGroupVersionKind(r.localGVK)
	current := object.DeepCopy()
	err := r.localMgr.GetClient().Get(r.GetContext(), types.NamespacedName{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}, current)
	if apierrors.IsNotFound(err) { // already deleted
		return nil
	}
	if err != nil {
		return err
	}

	ownerClusterID := current.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]

	if ownerClusterID == "" {
		log.Info("deletion is skipped, object is owned by this cluster")

		return nil
	}

	if r.isOwnedByAnotherAliveCluster(ownerClusterID) {
		log.Info("deletion is skipped, owned by another live cluster")

		return nil
	}

	err = r.localMgr.GetClient().Delete(r.GetContext(), object)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	return nil
}

func (r *syncReconciler) isOwnedByAnotherAliveCluster(ownerClusterID string) bool {
	return ownerClusterID != "" && r.clustersManager.GetAliveClustersByID()[ownerClusterID] != nil && ownerClusterID != r.clusterID
}

func (r *syncReconciler) getObjectDesiredState() *util.DynamicDesiredState {
	return &util.DynamicDesiredState{
		ShouldCreateFunc: func(desired runtime.Object) (bool, error) {
			metaObj, err := meta.Accessor(desired)
			if err != nil {
				return false, err
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

func (r *syncReconciler) initLocalInformer(ctx context.Context, obj *unstructured.Unstructured) error {
	key := obj.GetObjectKind().GroupVersionKind().String()

	if _, ok := r.localInformers[key]; ok {
		return nil
	}

	r.GetLogger().Info("init local informer", "gvk", key)

	localInformer, err := r.localMgr.GetCache().GetInformer(ctx, obj)
	if err != nil {
		return errors.WrapIf(err, "could not create local informer for clusters")
	}

	err = r.ctrl.Watch(&source.Informer{
		Informer: localInformer,
	}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(obj handler.MapObject) []reconcile.Request {
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Name:      obj.Meta.GetName(),
						Namespace: obj.Meta.GetNamespace(),
					},
				},
			}
		}),
	}, r.localPredicate())
	if err != nil {
		return errors.WrapIf(err, "could not create watch for local informer")
	}

	r.localInformers[key] = struct{}{}

	return nil
}

func (r *syncReconciler) localPredicate() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			_, ok := e.Meta.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]

			return ok
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			_, ok := e.MetaOld.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]

			return ok
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			_, ok := e.Meta.GetAnnotations()[clusterregistryv1alpha1.OwnershipAnnotation]

			return ok
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
}
