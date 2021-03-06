// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	"github.com/go-logr/logr"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcilerPolkadot) handleService(CRInstance *polkadotv1alpha1.Polkadot) (bool, error) {
	handler := getHandlerService(CRInstance)
	return handler.handleServiceSpecific(r,CRInstance)
}

//pattern factory
func getHandlerService(CRInstance *polkadotv1alpha1.Polkadot) IHandlerService {
	if CRKind(CRInstance.Spec.Kind) == Validator {
		return &handlerServiceValidator{}
	}
	if CRKind(CRInstance.Spec.Kind) == Sentry {
		return &handlerServiceSentry{}
	}
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		return &handlerServiceSentryAndValidator{}
	}
	return &handlerServiceDefault{}
}

//pattern Strategy
type IHandlerService interface {
	handleServiceSpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error)
}

type handlerServiceValidator struct {
}
func (h *handlerServiceValidator) handleServiceSpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error) {
	return r.handleServiceGeneric(CRInstance, newServiceValidator(CRInstance))
}

type handlerServiceSentry struct {
}
func (h *handlerServiceSentry) handleServiceSpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error) {
	return r.handleServiceGeneric(CRInstance, newServiceSentry(CRInstance))
}

type handlerServiceSentryAndValidator struct {
}
func (h *handlerServiceSentryAndValidator) handleServiceSpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error) {
	isForcedRequeue, err := r.handleServiceGeneric(CRInstance, newServiceSentry(CRInstance))
	if isForcedRequeue == ForcedRequeue || err != nil {
		return isForcedRequeue, err
	}
	return r.handleServiceGeneric(CRInstance, newServiceValidator(CRInstance))
}

type handlerServiceDefault struct {
}
func (h *handlerServiceDefault) handleServiceSpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error){
	return handleSkip()
}

func (r *ReconcilerPolkadot) handleServiceGeneric(CRInstance *polkadotv1alpha1.Polkadot, desiredResource *corev1.Service) (bool, error) {

	logger := log.WithValues("Service.Namespace", desiredResource.Namespace, "Service.Name", desiredResource.Name)

	toBeFoundResource := &corev1.Service{}
	isNotFound,err := r.fetchResource(toBeFoundResource,types.NamespacedName{Name: desiredResource.Name, Namespace: desiredResource.Namespace})
	if err != nil {
		logger.Error(err, "Error on fetch the Service...")
		return NotForcedRequeue, err
	}
	if isNotFound == true {
		logger.Info("Service not found...")
		logger.Info("Creating a new Service...")
		err := r.createResource(desiredResource, CRInstance)
		if err != nil {
			logger.Error(err, "Error on creating a new Service...")
			return NotForcedRequeue, err
		}
		logger.Info("Created the new Service")
		return ForcedRequeue, nil
	}
	foundResource := toBeFoundResource

	if areServicesDifferent(foundResource, desiredResource, logger) {
		logger.Info("Updating the Service...")
		err := r.updateResource(desiredResource)
		if err != nil {
			logger.Error(err, "Update Service Error...")
			return NotForcedRequeue, err
		}
		logger.Info("Updated the Service...")
	}

	return NotForcedRequeue, nil
}

func areServicesDifferent(currentService *corev1.Service, desiredService *corev1.Service, logger logr.Logger) bool {
	result := false
	return result
}