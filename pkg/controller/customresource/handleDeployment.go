package customresource

import (
	"context"
	"github.com/go-logr/logr"
	cachev1alpha1 "github.com/ironoa/kubernetes-customresource-operator/pkg/apis/cache/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileCustomResource) handleDeployment(CRInstance *cachev1alpha1.CustomResource) (bool, error) {
	const NotForcedRequeue = false
	const ForcedRequeue = true

	desiredDeployment := newDeploymentForCR(CRInstance)
	logger := log.WithValues("Deployment.Namespace", desiredDeployment.Namespace, "Deployment.Name", desiredDeployment.Name)

	foundDeployment, err := r.fetchDeployment(desiredDeployment)
	if err != nil {
		logger.Error(err, "Error on fetch the Deployment...")
		return NotForcedRequeue, err
	}
	if foundDeployment == nil {
		logger.Info("Deployment not found...")
		logger.Info("Creating a new Deployment...")
		err := r.createDeployment(desiredDeployment, CRInstance, logger)
		if err != nil {
			logger.Error(err, "Error on creating a new Deployment...")
			return NotForcedRequeue, err
		}
		logger.Info("Created the new Deployment")
		return ForcedRequeue, nil
	}

	if areDeploymentsDifferent(foundDeployment, desiredDeployment, logger) {
		err := r.updateDeployment(desiredDeployment, logger)
		if err != nil {
			logger.Error(err, "Update Deployment Error...")
			return NotForcedRequeue, err
		}
		logger.Info("Updated the Deployment...")
	}

	return NotForcedRequeue, nil
}

func (r *ReconcileCustomResource) fetchDeployment(deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return nil, nil
	}
	return found, err
}

func (r *ReconcileCustomResource) createDeployment(deployment *appsv1.Deployment, CRInstance *cachev1alpha1.CustomResource, logger logr.Logger) error {
	err := r.setOwnership(CRInstance, deployment)
	if err != nil {
		logger.Error(err, "Error on setting the ownership...")
		return err
	}
	err = r.client.Create(context.TODO(), deployment)
	return err
}

func areDeploymentsDifferent(currentDeployment *appsv1.Deployment, desiredDeployment *appsv1.Deployment, logger logr.Logger) bool {
	result := false

	if isDeploymentReplicaDifferent(currentDeployment, desiredDeployment, logger) {
		result = true
	}
	if isDeploymentVersionDifferent(currentDeployment, desiredDeployment, logger) {
		result = true
	}

	return result
}

func isDeploymentReplicaDifferent(currentDeployment *appsv1.Deployment, desiredDeployment *appsv1.Deployment, logger logr.Logger) bool {
	size := *desiredDeployment.Spec.Replicas
	if *currentDeployment.Spec.Replicas != size {
		logger.Info("Find a replica size mismatch...")
		return true
	}
	return false
}

func isDeploymentVersionDifferent(currentDeployment *appsv1.Deployment, desiredDeployment *appsv1.Deployment, logger logr.Logger) bool {
	version := desiredDeployment.ObjectMeta.Labels["version"]
	if currentDeployment.ObjectMeta.Labels["version"] != version {
		logger.Info("Found a version mismatch...")
		return true
	}
	return false
}

func (r *ReconcileCustomResource) updateDeployment(deployment *appsv1.Deployment, logger logr.Logger) error {
	logger.Info("Updating the Deployment...")
	return r.client.Update(context.TODO(), deployment)
}
