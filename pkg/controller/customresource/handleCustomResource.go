package customresource

import (
	"context"
	cachev1alpha1 "github.com/ironoa/kubernetes-customresource-operator/pkg/apis/cache/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileCustomResource) handleCustomResource(request reconcile.Request) (*cachev1alpha1.CustomResource, error) {
	logger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	found, err := r.fetchCustomResource(request)
	if err != nil {
		logger.Error(err, "Error on fetch the Custom Resource...")
		return nil, err
	}
	if found == nil {
		logger.Info("Custom Resource not found...")
		return nil, nil
	}

	return found, nil
}

func (r *ReconcileCustomResource) fetchCustomResource(request reconcile.Request) (*cachev1alpha1.CustomResource, error) {
	found := &cachev1alpha1.CustomResource{}
	err := r.client.Get(context.TODO(), request.NamespacedName, found)
	if err != nil && errors.IsNotFound(err) {
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		// Return and don't requeue
		return nil, nil
	}
	return found, err
}