// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcilerPolkadot) handleCustomResource(request reconcile.Request) (*polkadotv1alpha1.Polkadot, error) {
	logger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	toBeFound := &polkadotv1alpha1.Polkadot{}
	isNotFound, err := r.fetchResource(toBeFound, types.NamespacedName{Name: request.Name, Namespace: request.Namespace})
	if err != nil {
		logger.Error(err, "Error on fetch the Custom Resource...")
		return nil, err
	}
	if isNotFound == true {
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		// Return and don't requeue
		logger.Info("Custom Resource not found...")
		return nil, nil
	}

	return toBeFound, nil
}