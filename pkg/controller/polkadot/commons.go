package polkadot

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcilePolkadot) setOwnership(owner metav1.Object, owned metav1.Object) error {
	return controllerutil.SetControllerReference(owner, owned, r.scheme)
}

type CRKind string
const (
	Sentry CRKind = "Sentry"
	Validator CRKind = "Validator"
	SentryAndValidator CRKind = "SentryAndValidator"
)

const(
	NotForcedRequeue = false
	ForcedRequeue = true
)

func handleSkip() (bool,error){
	return NotForcedRequeue,nil
}