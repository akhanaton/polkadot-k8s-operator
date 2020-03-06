// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	"github.com/go-logr/logr"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_polkadot")

// ReconcilerPolkadot reconciles a Polkadot object
type ReconcilerPolkadot struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Add creates a new Polkadot Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilerPolkadot{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("polkadot-controller", mgr, controller.Options{Reconciler: r, MaxConcurrentReconciles: 1})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CustomResource
	err = c.Watch(&source.Kind{Type: &polkadotv1alpha1.Polkadot{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource StatefulSet and requeue the owner CustomResource
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &polkadotv1alpha1.Polkadot{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Service and requeue the owner CustomResource
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &polkadotv1alpha1.Polkadot{},
	})
	if err != nil {
		return err
	}

	//TODO add watch for NetworkPolicy

	return nil
}

// blank assignment to verify that ReconcilerPolkadot implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilerPolkadot{}

// Reconcile reads that state of the cluster for a CustomResource object and makes changes based on the state read
// and what is in the CustomResource.Spec
func (r *ReconcilerPolkadot) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	logger.Info("Reconciling CustomResource")

	handledCRInstance, err := r.handleCustomResource(request)
	if err != nil {
		return handleRequeueError(err,logger)
	}
	if handledCRInstance == nil {
		return handleRequeueStd(err, logger)
	}

	isRequeueForced, err := r.handleStatefulSet(handledCRInstance)
	if err != nil {
		return handleRequeueError(err,logger)
	}
	if isRequeueForced {
		return handleRequeueForced(err, logger)
	}

	isRequeueForced, err = r.handleService(handledCRInstance)
	if err != nil {
		return handleRequeueError(err,logger)
	}
	if isRequeueForced {
		return handleRequeueForced(err, logger)
	}

	isRequeueForced, err = r.handleNetworkPolicy(handledCRInstance)
	if err != nil {
		return handleRequeueError(err,logger)
	}
	if isRequeueForced {
		return handleRequeueForced(err, logger)
	}

	return handleRequeueStd(err, logger)
}

// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func handleRequeueError (err error, logger logr.Logger) (reconcile.Result, error){
	logger.Info("Requeing the Reconciling request... ")
	return reconcile.Result{}, err
}

func handleRequeueForced (err error, logger logr.Logger) (reconcile.Result, error){
	logger.Info("Requeing the Reconciling request... ")
	return reconcile.Result{Requeue: true}, nil
}

func handleRequeueStd (err error, logger logr.Logger) (reconcile.Result, error){
	logger.Info("Return and not requeing the request")
	return reconcile.Result{}, nil
}