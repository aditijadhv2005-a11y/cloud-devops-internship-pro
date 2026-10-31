package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagedAppSpec describes the desired application workload.
type ManagedAppSpec struct {
	// Image is an immutable tag or digest recommended for reproducible deployments.
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`
	// Replicas is the desired number of application replicas.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=20
	// +kubebuilder:default=2
	Replicas *int32 `json:"replicas,omitempty"`
	// Port is the container and Service target port.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=8080
	Port int32 `json:"port,omitempty"`
	// Resources is the container resource budget.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type ManagedAppStatus struct {
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=mapp
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image"
// +kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.readyReplicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ManagedApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ManagedAppSpec   `json:"spec,omitempty"`
	Status            ManagedAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ManagedAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedApp `json:"items"`
}

func init() { SchemeBuilder.Register(&ManagedApp{}, &ManagedAppList{}) }
