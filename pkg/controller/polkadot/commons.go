// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	"context"
	"github.com/go-logr/logr"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type CRKind string
const (
	Sentry CRKind = "Sentry"
	Validator CRKind = "Validator"
	SentryAndValidator CRKind = "SentryAndValidator"
)

const(
	NotForcedRequeue = false
	ForcedRequeue = true
)

func handleSkip() (bool,error){
	return NotForcedRequeue,nil
}

func (r *ReconcilerPolkadot) setOwnership(owner metav1.Object, owned metav1.Object) error {
	return controllerutil.SetControllerReference(owner, owned, r.scheme)
}

func (r *ReconcilerPolkadot) createResource(resource interface{}, CRInstance *polkadotv1alpha1.Polkadot, logger logr.Logger) error {
	metaObj := resource.(metav1.Object)
	err := r.setOwnership(CRInstance, metaObj)
	if err != nil {
		logger.Error(err, "Error on setting the ownership...")
		return err
	}
	runObj := resource.(runtime.Object)
	return r.client.Create(context.TODO(), runObj)
}