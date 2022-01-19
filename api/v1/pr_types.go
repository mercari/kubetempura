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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PRSpec defines the desired state of PR
type PRSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// The parent review app name
	ParentReviewApp string `json:"parentReviewApp"`

	// +kubebuilder:validation:Required
	// PR Number
	PRNumber string `json:"prNumber"`

	// +kubebuilder:validation:Required
	// The sha of the latest commit.
	HeadCommitRef string `json:"headCommitRef"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// Environment variables for adding / overriding the default values.
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`
}

// PRStatus defines the observed state of PR
type PRStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PR is the Schema for the prs API.
// PR is the internal CRD for each PRs. The GitHub Webhooks' handler will create/update/delete this resource.
type PR struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PRSpec   `json:"spec,omitempty"`
	Status PRStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PRList contains a list of PR
type PRList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PR `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PR{}, &PRList{})
}
