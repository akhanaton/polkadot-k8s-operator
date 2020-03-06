// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	"context"
	"github.com/go-logr/logr"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcilerPolkadot) handleStatefulSet(CRInstance *polkadotv1alpha1.Polkadot) (bool, error){
	handler := getHandlerStatefulSet(CRInstance)
	return handler.handleStatefulSetSpecific(r,CRInstance)
}

//pattern factory
func getHandlerStatefulSet(CRInstance *polkadotv1alpha1.Polkadot) IHandlerStatefulSet {
	if CRKind(CRInstance.Spec.Kind) == Validator {
		return &handlerStatefulSetValidator{}
	}
	if CRKind(CRInstance.Spec.Kind) == Sentry {
		return &handlerStatefulSetSentry{}
	}
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		return &handlerStatefulSetSentryAndValidator{}
	}
	return &handlerStatefulSetDefault{}
}

//pattern Strategy
type IHandlerStatefulSet interface {
	handleStatefulSetSpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error)
}

type handlerStatefulSetValidator struct {
}
func (h *handlerStatefulSetValidator) handleStatefulSetSpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error){
	return r.handleStatefulSetGeneric(CRInstance, newValidatorStatefulSetForCR(CRInstance))
}

type handlerStatefulSetSentry struct {
}
func (h *handlerStatefulSetSentry) handleStatefulSetSpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error){
	return r.handleStatefulSetGeneric(CRInstance, newSentryStatefulSetForCR(CRInstance))
}

type handlerStatefulSetSentryAndValidator struct {
}
func (h *handlerStatefulSetSentryAndValidator) handleStatefulSetSpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error){
	isForcedRequeue, err := r.handleStatefulSetGeneric(CRInstance, newSentryStatefulSetForCR(CRInstance))
	if isForcedRequeue == ForcedRequeue || err != nil {
		return isForcedRequeue, err
	}
	return r.handleStatefulSetGeneric(CRInstance, newValidatorStatefulSetForCR(CRInstance))
}

type handlerStatefulSetDefault struct {
}
func (h *handlerStatefulSetDefault) handleStatefulSetSpecific(r *ReconcilerPolkadot, CRInstance *polkadotv1alpha1.Polkadot) (bool, error){
	return handleSkip()
}

func (r *ReconcilerPolkadot) handleStatefulSetGeneric(CRInstance *polkadotv1alpha1.Polkadot, desiredResource *appsv1.StatefulSet) (bool, error) {

	logger := log.WithValues("Deployment.Namespace", desiredResource.Namespace, "Deployment.Name", desiredResource.Name)

	toBeFound := &appsv1.StatefulSet{}
	isNotFound, err := r.fetchResource(toBeFound,types.NamespacedName{Name: desiredResource.Name, Namespace: desiredResource.Namespace})
	if err != nil {
		logger.Error(err, "Error on fetch the StatefulSet...")
		return NotForcedRequeue, err
	}
	if isNotFound == true {
		logger.Info("StatefulSet not found...")
		logger.Info("Creating a new StatefulSet...")
		err := r.createResource(desiredResource, CRInstance, logger)
		if err != nil {
			logger.Error(err, "Error on creating a new StatefulSet...")
			return NotForcedRequeue, err
		}
		logger.Info("Created the new StatefulSet")
		return ForcedRequeue, nil
	}
	foundResource := toBeFound

	if areStatefulSetDifferent(foundResource, desiredResource, logger) {
		logger.Info("Updating the StatefulSet...")
		err := r.updateStatefulSet(desiredResource)
		if err != nil {
			logger.Error(err, "Update StatefulSet Error...")
			return NotForcedRequeue, err
		}
		logger.Info("Updated the StatefulSet...")
	}

	return NotForcedRequeue, nil
}

func (r *ReconcilerPolkadot) updateStatefulSet(obj *appsv1.StatefulSet) error {
	return r.client.Update(context.TODO(), obj)
}

func areStatefulSetDifferent(current *appsv1.StatefulSet, desired *appsv1.StatefulSet, logger logr.Logger) bool {
	result := false

	if isStatefulSetReplicaDifferent(current, desired, logger) {
		result = true
	}
	if isStatefulSetVersionDifferent(current, desired, logger) {
		result = true
	}

	return result
}

func isStatefulSetReplicaDifferent(current *appsv1.StatefulSet, desired *appsv1.StatefulSet, logger logr.Logger) bool {
	size := *desired.Spec.Replicas
	if *current.Spec.Replicas != size {
		logger.Info("Found a replica size mismatch...")
		return true
	}
	return false
}

func isStatefulSetVersionDifferent(current *appsv1.StatefulSet, desired *appsv1.StatefulSet, logger logr.Logger) bool {
	version := desired.ObjectMeta.Labels["version"]
	if current.ObjectMeta.Labels["version"] != version {
		logger.Info("Found a version mismatch...")
		return true
	}
	return false
}
