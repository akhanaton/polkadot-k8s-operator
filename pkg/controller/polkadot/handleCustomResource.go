package polkadot

import (
	"context"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcilePolkadot) handleCustomResource(request reconcile.Request) (*polkadotv1alpha1.Polkadot, error) {
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

func (r *ReconcilePolkadot) fetchCustomResource(request reconcile.Request) (*polkadotv1alpha1.Polkadot, error) {
	found := &polkadotv1alpha1.Polkadot{}
	err := r.client.Get(context.TODO(), request.NamespacedName, found)
	if err != nil && errors.IsNotFound(err) {
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		// Return and don't requeue
		return nil, nil
	}
	return found, err
}