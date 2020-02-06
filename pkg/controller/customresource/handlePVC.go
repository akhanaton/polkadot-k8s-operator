package customresource

import (
	"context"
	"github.com/go-logr/logr"
	cachev1alpha1 "github.com/ironoa/kubernetes-customresource-operator/pkg/apis/cache/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileCustomResource) handlePVC(CRInstance *cachev1alpha1.CustomResource) (bool, error) {
	const NotForcedRequeue = false
	const ForcedRequeue = true

	desiredPVC := newPVCForCR(CRInstance)
	logger := log.WithValues("PVC.Namespace", desiredPVC.Namespace, "PVC.Name", desiredPVC.Name)

	foundPVC, err := r.fetchPVC(desiredPVC)
	if err != nil {
		logger.Error(err, "Error on fetch the PVC...")
		return NotForcedRequeue, err
	}
	if foundPVC == nil {
		logger.Info("PVC not found...")
		logger.Info("Creating a new PVC...")
		err := r.createPVC(desiredPVC, CRInstance, logger)
		if err != nil {
			logger.Error(err, "Error on creating a new PVC...")
			return NotForcedRequeue, err
		}
		logger.Info("Created the new PVC")
		return ForcedRequeue, nil
	}

	if arePVCsDifferent(foundPVC, desiredPVC, logger) {
		err := r.updatePVC(desiredPVC, logger)
		if err != nil {
			logger.Error(err, "Update PVC Error...")
			return NotForcedRequeue, err
		}
		logger.Info("Updated the PVC...")
	}

	return NotForcedRequeue, nil
}

func (r *ReconcileCustomResource) fetchPVC(PVC *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	found := &corev1.PersistentVolumeClaim{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: PVC.Name, Namespace: PVC.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return nil, nil
	}
	return found, err
}

func (r *ReconcileCustomResource) createPVC(PVC *corev1.PersistentVolumeClaim, CRInstance *cachev1alpha1.CustomResource, logger logr.Logger) error {
	err := r.setOwnership(CRInstance, PVC)
	if err != nil {
		logger.Error(err, "Error on setting the ownership...")
		return err
	}
	return r.client.Create(context.TODO(), PVC)
}

func arePVCsDifferent(currentPVC *corev1.PersistentVolumeClaim, desiredPVC *corev1.PersistentVolumeClaim, logger logr.Logger) bool {
	result := false
	return result
}

func (r *ReconcileCustomResource) updatePVC(PVC *corev1.PersistentVolumeClaim, logger logr.Logger) error {
	logger.Info("Updating the Persistent Volume Claim...")
	return r.client.Update(context.TODO(), PVC)
}
