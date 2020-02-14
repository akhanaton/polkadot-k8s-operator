package customresource

import (
	"context"
	"github.com/go-logr/logr"
	cachev1alpha1 "github.com/ironoa/kubernetes-customresource-operator/pkg/apis/cache/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileCustomResource) handleService(CRInstance *cachev1alpha1.CustomResource) (bool, error) {
	const NotForcedRequeue = false
	const ForcedRequeue = true

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

	return NotForcedRequeue, nil // TODO handle default
}

func (r *ReconcileCustomResource) handleSpecificService(CRInstance *cachev1alpha1.CustomResource, desiredService *corev1.Service) (bool, error) {
	const NotForcedRequeue = false
	const ForcedRequeue = true

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

func (r *ReconcileCustomResource) fetchService(service *corev1.Service) (*corev1.Service, error) {
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return nil, nil
	}
	return found, err
}

func (r *ReconcileCustomResource) createService(service *corev1.Service, CRInstance *cachev1alpha1.CustomResource, logger logr.Logger) error {
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

func (r *ReconcileCustomResource) updateService(service *corev1.Service, logger logr.Logger) error {
	return r.client.Update(context.TODO(), service)
}