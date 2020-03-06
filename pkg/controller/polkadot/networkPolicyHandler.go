// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	v1 "k8s.io/api/networking/v1"
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
	return r.handleNetworkPolicyGeneric(CRInstance, newNetworkPolicyValidator(CRInstance))
}

type handlerNetworkPolicyDefault struct {
}
func (h *handlerNetworkPolicyDefault) handleNetworkPolicySpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error){
	return handleSkip()
}

func (r *ReconcilerPolkadot) handleNetworkPolicyGeneric(CRInstance *polkadotv1alpha1.Polkadot, desiredResource *v1.NetworkPolicy) (bool, error) {

	logger := log.WithValues("Service.Namespace", desiredResource.Namespace, "Service.Name", desiredResource.Name)

	toBeFoundResource := &v1.NetworkPolicy{}
	isNotFound,err := r.fetchResource(toBeFoundResource,types.NamespacedName{Name: desiredResource.Name, Namespace: desiredResource.Namespace})
	if err != nil {
		logger.Error(err, "Error on fetch the Network Policy...")
		return NotForcedRequeue, err
	}
	if isNotFound == true {
		logger.Info("Network Policy not found...")
		logger.Info("Creating a new Network Policy...")
		err := r.createResource(desiredResource, CRInstance, logger)
		if err != nil {
			logger.Error(err, "Error on creating a new Network Policy...")
			return NotForcedRequeue, err
		}
		logger.Info("Created the new Network Policy")
		return ForcedRequeue, nil
	}

	//TODO add check differences

	return NotForcedRequeue, nil
}