/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IdlingResourceSpec defines the desired state of IdlingResource
type IdlingResourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The reference to the idle-able resource
	IdlingResourceRef CrossVersionObjectReference `json:"idlingResourceRef"`

	// The desired state of idling. Defaults to false.
	// +kubebuilder:default:false
	Idle bool `json:"idle"`

	// The number of replicas after the wakeup phase
	// +optional
	// +kubebuilder:validation:Minimum=1
	ResumeReplicas *int32 `json:"resumeReplicas,omitempty"`
}

// CrossVersionObjectReference contains enough information to let you identify the referred resource.
type CrossVersionObjectReference struct {
	// Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds"
	Kind string `json:"kind"`

	// Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string `json:"name"`

	// API version of the referent
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`
}

// IdlingResourceStatus defines the observed state of IdlingResource
type IdlingResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	PreviousReplicas *int32 `json:"previousReplicas"`
}

// +kubebuilder:resource:shortName=ir
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Idle",type="boolean",JSONPath=".spec.idle"
// +kubebuilder:printcolumn:name="RefKind",type="string",JSONPath=".spec.idlingResourceRef.kind"
// +kubebuilder:printcolumn:name="RefName",type="string",JSONPath=".spec.idlingResourceRef.name"

// IdlingResource is the Schema for the idlingresources API
type IdlingResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IdlingResourceSpec   `json:"spec,omitempty"`
	Status IdlingResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IdlingResourceList contains a list of IdlingResource
type IdlingResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IdlingResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IdlingResource{}, &IdlingResourceList{})
}
