package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CustomResourceSpec defines the desired state of CustomResource
type CustomResourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	ClientVersion string `json:"clientVersion"`
	Kind string `json:"kind"`
	Validator `json:"validator,omitempty"`
	Sentry `json:"sentry,omitempty"`
}

type Validator struct {
	ClientName string `json:"clientName"`
	NodeKey string `json:"nodeKey"`
	ReservedSentryID string `json:"reservedSentryID,omitempty"`
	CPULimit string `json:"CPULimit,omitempty"`
	MemoryLimit string `json:"memoryLimit,omitempty"`
}

type Sentry struct {
	Replicas int32 `json:"replicas"`
	ClientName string `json:"clientName"`
	NodeKey string `json:"nodeKey"`
	ReservedValidatorID string `json:"reservedValidatorID,omitempty"`
	CPULimit string `json:"CPULimit,omitempty"`
	MemoryLimit string `json:"memoryLimit,omitempty"`
}

// CustomResourceStatus defines the observed state of CustomResource
type CustomResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Nodes are the names of the CustomResource pods... ?? to check
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomResource is the Schema for the customresources API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=customresources,scope=Namespaced
type CustomResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomResourceSpec   `json:"spec,omitempty"`
	Status CustomResourceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomResourceList contains a list of CustomResource
type CustomResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomResource{}, &CustomResourceList{})
}
