package customresource

import (
	"context"
	"github.com/go-logr/logr"
	cachev1alpha1 "github.com/ironoa/kubernetes-customresource-operator/pkg/apis/cache/v1alpha1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileCustomResource) handleNetworkPolicy(CRInstance *cachev1alpha1.CustomResource) (bool, error) {

	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		return r.handleSpecificNetworkPolicy(CRInstance, newValidatorNetworkPolicyForCR(CRInstance))
	}

	return defaultHandler()
}

func (r *ReconcileCustomResource) handleSpecificNetworkPolicy(CRInstance *cachev1alpha1.CustomResource, desiredNetworkPolicy *v1.NetworkPolicy) (bool, error) {

	logger := log.WithValues("Service.Namespace", desiredNetworkPolicy.Namespace, "Service.Name", desiredNetworkPolicy.Name)

	foundNP, err := r.fetchNP(desiredNetworkPolicy)
	if err != nil {
		logger.Error(err, "Error on fetch the Network Policy...")
		return NotForcedRequeue, err
	}
	if foundNP == nil {
		logger.Info("Network Policy not found...")
		logger.Info("Creating a new Network Policy...")
		err := r.createNP(desiredNetworkPolicy, CRInstance, logger)
		if err != nil {
			logger.Error(err, "Error on creating a new Network Policy...")
			return NotForcedRequeue, err
		}
		logger.Info("Created the new Network Policy")
		return ForcedRequeue, nil
	}

	return NotForcedRequeue, nil
}

func (r *ReconcileCustomResource) fetchNP(np *v1.NetworkPolicy) (*v1.NetworkPolicy, error) {
	found := &v1.NetworkPolicy{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: np.Name, Namespace: np.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return nil, nil
	}
	return found, err
}

func (r *ReconcileCustomResource) createNP(networkPolicy *v1.NetworkPolicy, CRInstance *cachev1alpha1.CustomResource, logger logr.Logger) error {
	err := r.setOwnership(CRInstance, networkPolicy)
	if err != nil {
		logger.Error(err, "Error on setting the ownership...")
		return err
	}
	return r.client.Create(context.TODO(), networkPolicy)
}