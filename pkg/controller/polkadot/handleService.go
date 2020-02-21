// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	"context"
	"github.com/go-logr/logr"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcilePolkadot) handleService(CRInstance *polkadotv1alpha1.Polkadot) (bool, error) {

	if CRKind(CRInstance.Spec.Kind) == Validator {
		return r.handleSpecificService(CRInstance, newValidatorServiceForCR(CRInstance))
	}
	if CRKind(CRInstance.Spec.Kind) == Sentry {
		return r.handleSpecificService(CRInstance, newSentryServiceForCR(CRInstance))
	}
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		isForcedRequeue, err := r.handleSpecificService(CRInstance, newSentryServiceForCR(CRInstance))
		if isForcedRequeue == ForcedRequeue || err != nil {
			return isForcedRequeue, err
		}
		return r.handleSpecificService(CRInstance, newValidatorServiceForCR(CRInstance))
	}

	return defaultHandler()
}

func (r *ReconcilePolkadot) handleSpecificService(CRInstance *polkadotv1alpha1.Polkadot, desiredService *corev1.Service) (bool, error) {

	logger := log.WithValues("Service.Namespace", desiredService.Namespace, "Service.Name", desiredService.Name)

	foundService, err := r.fetchService(desiredService)
	if err != nil {
		logger.Error(err, "Error on fetch the Service...")
		return NotForcedRequeue, err
	}
	if foundService == nil {
		logger.Info("Service not found...")
		logger.Info("Creating a new Service...")
		err := r.createService(desiredService, CRInstance, logger)
		if err != nil {
			logger.Error(err, "Error on creating a new Service...")
			return NotForcedRequeue, err
		}
		logger.Info("Created the new Service")
		return ForcedRequeue, nil
	}

	if areServicesDifferent(foundService, desiredService, logger) {
		logger.Info("Updating the Service...")
		err := r.updateService(desiredService, logger)
		if err != nil {
			logger.Error(err, "Update Service Error...")
			return NotForcedRequeue, err
		}
		logger.Info("Updated the Service...")
	}

	return NotForcedRequeue, nil
}

func (r *ReconcilePolkadot) fetchService(service *corev1.Service) (*corev1.Service, error) {
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return nil, nil
	}
	return found, err
}

func (r *ReconcilePolkadot) createService(service *corev1.Service, CRInstance *polkadotv1alpha1.Polkadot, logger logr.Logger) error {
	err := r.setOwnership(CRInstance, service)
	if err != nil {
		logger.Error(err, "Error on setting the ownership...")
		return err
	}
	return r.client.Create(context.TODO(), service)
}

func areServicesDifferent(currentService *corev1.Service, desiredService *corev1.Service, logger logr.Logger) bool {
	result := false
	return result
}

func (r *ReconcilePolkadot) updateService(service *corev1.Service, logger logr.Logger) error {
	return r.client.Update(context.TODO(), service)
}