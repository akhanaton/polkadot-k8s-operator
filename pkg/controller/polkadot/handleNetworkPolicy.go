// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	"context"
	"github.com/go-logr/logr"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcilerPolkadot) handleNetworkPolicy(CRInstance *polkadotv1alpha1.Polkadot) (bool, error) {
	handler := getHandlerNetworkPolicy(CRInstance)
	return handler.handleNetworkPolicySpecific(r,CRInstance)
}

//pattern factory
func getHandlerNetworkPolicy(CRInstance *polkadotv1alpha1.Polkadot) IHandlerNetworkPolicy {
	if CRInstance.Spec.IsNetworkPolicyActive != "true" {
		return &handlerNetworkPolicyDefault{}
	}
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		return &handlerNetworkPolicySentryAndValidator{}
	}
	return &handlerNetworkPolicyDefault{}
}

//pattern Strategy
type IHandlerNetworkPolicy interface {
	handleNetworkPolicySpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error)
}

type handlerNetworkPolicySentryAndValidator struct {
}
func (h *handlerNetworkPolicySentryAndValidator) handleNetworkPolicySpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error) {
	return r.handleNetworkPolicyGeneric(CRInstance,newValidatorNetworkPolicyForCR(CRInstance))
}

type handlerNetworkPolicyDefault struct {
}
func (h *handlerNetworkPolicyDefault) handleNetworkPolicySpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error){
	return handleSkip()
}

func (r *ReconcilerPolkadot) handleNetworkPolicyGeneric(CRInstance *polkadotv1alpha1.Polkadot, desiredNetworkPolicy *v1.NetworkPolicy) (bool, error) {

	logger := log.WithValues("Service.Namespace", desiredNetworkPolicy.Namespace, "Service.Name", desiredNetworkPolicy.Name)

	foundNP, err := r.fetchNetworkPolicy(desiredNetworkPolicy)
	if err != nil {
		logger.Error(err, "Error on fetch the Network Policy...")
		return NotForcedRequeue, err
	}
	if foundNP == nil {
		logger.Info("Network Policy not found...")
		logger.Info("Creating a new Network Policy...")
		err := r.createNetworkPolicy(desiredNetworkPolicy, CRInstance, logger)
		if err != nil {
			logger.Error(err, "Error on creating a new Network Policy...")
			return NotForcedRequeue, err
		}
		logger.Info("Created the new Network Policy")
		return ForcedRequeue, nil
	}

	return NotForcedRequeue, nil
}

func (r *ReconcilerPolkadot) fetchNetworkPolicy(np *v1.NetworkPolicy) (*v1.NetworkPolicy, error) {
	found := &v1.NetworkPolicy{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: np.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return nil, nil
	}
	return found, err
}

func (r *ReconcilerPolkadot) createNetworkPolicy(networkPolicy *v1.NetworkPolicy, CRInstance *polkadotv1alpha1.Polkadot, logger logr.Logger) error {
	err := r.setOwnership(CRInstance, networkPolicy)
	if err != nil {
		logger.Error(err, "Error on setting the ownership...")
		return err
	}
	return r.client.Create(context.TODO(), networkPolicy)
}