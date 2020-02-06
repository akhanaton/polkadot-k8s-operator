package customresource

import (
	"context"
	"github.com/go-logr/logr"
	cachev1alpha1 "github.com/ironoa/kubernetes-customresource-operator/pkg/apis/cache/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileCustomResource) handleStatefulSet(CRInstance *cachev1alpha1.CustomResource, desiredResource *appsv1.StatefulSet) (bool, error) {
	const NotForcedRequeue = false
	const ForcedRequeue = true

	logger := log.WithValues("Deployment.Namespace", desiredResource.Namespace, "Deployment.Name", desiredResource.Name)

	foundResource, err := r.fetchStatefulSet(desiredResource)
	if err != nil {
		logger.Error(err, "Error on fetch the StatefulSet...")
		return NotForcedRequeue, err
	}
	if foundResource == nil {
		logger.Info("StatefulSet not found...")
		logger.Info("Creating a new StatefulSet...")
		err := r.createStatefulSet(desiredResource, CRInstance, logger)
		if err != nil {
			logger.Error(err, "Error on creating a new StatefulSet...")
			return NotForcedRequeue, err
		}
		logger.Info("Created the new StatefulSet")
		return ForcedRequeue, nil
	}

	if areStatefulSetDifferent(foundResource, desiredResource, logger) {
		logger.Info("Updating the StatefulSet...")
		err := r.updateStatefulSet(desiredResource, logger)
		if err != nil {
			logger.Error(err, "Update StatefulSet Error...")
			return NotForcedRequeue, err
		}
		logger.Info("Updated the StatefulSet...")
	}

	return NotForcedRequeue, nil
}

func (r *ReconcileCustomResource) fetchStatefulSet(obj *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	found := &appsv1.StatefulSet{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return nil, nil
	}
	return found, err
}

func (r *ReconcileCustomResource) createStatefulSet(statefulSet *appsv1.StatefulSet, CRInstance *cachev1alpha1.CustomResource, logger logr.Logger) error {
	err := r.setOwnership(CRInstance, statefulSet)
	if err != nil {
		logger.Error(err, "Error on setting the ownership...")
		return err
	}
	err = r.client.Create(context.TODO(), statefulSet)
	return err
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

func (r *ReconcileCustomResource) updateStatefulSet(obj *appsv1.StatefulSet, logger logr.Logger) error {
	return r.client.Update(context.TODO(), obj)
}
