package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PolkadotSpec defines the desired state of Polkadot
type PolkadotSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	ClientVersion string `json:"clientVersion"`
	Kind string `json:"kind"`
	Validator `json:"validator,omitempty"`
	Sentry `json:"sentry,omitempty"`
	IsNetworkPolicyActive string `json:"isNetworkPolicyActive,omitempty"`
	IsDataPersistenceActive string `json:"isDataPersistenceActive,omitempty"`
	IsMetricsSupportActive string `json:"isMetricsSupportActive,omitempty"`
}

type Validator struct {
	ClientName string `json:"clientName"`
	NodeKey string `json:"nodeKey"`
	ReservedSentryID string `json:"reservedSentryID,omitempty"`
	CPULimit string `json:"CPULimit,omitempty"`
	MemoryLimit string `json:"memoryLimit,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
}

type Sentry struct {
	Replicas int32 `json:"replicas"`
	ClientName string `json:"clientName"`
	NodeKey string `json:"nodeKey"`
	ReservedValidatorID string `json:"reservedValidatorID,omitempty"`
	CPULimit string `json:"CPULimit,omitempty"`
	MemoryLimit string `json:"memoryLimit,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
}

// PolkadotStatus defines the observed state of Polkadot
type PolkadotStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Nodes are the names of the CustomResource pods... ?? to check
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Polkadot is the Schema for the polkadots API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=polkadots,scope=Namespaced
type Polkadot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PolkadotSpec   `json:"spec,omitempty"`
	Status PolkadotStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolkadotList contains a list of Polkadot
type PolkadotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Polkadot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Polkadot{}, &PolkadotList{})
}
