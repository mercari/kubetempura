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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ReviewAppSpec defines the desired state of ReviewApp
type ReviewAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// The GitHub URL of the repository. E.g. https://github.com/kouzoh/mercari-echo-us
	GithubRepository string `json:"githubRepository"`

	// +kubebuilder:validation:Required
	// +kubebuilder:pruning:PreserveUnknownFields
	// Resource is a field for any kind of Kubernetes resources. It can be deployment or service or anything.
	Resources []unstructured.Unstructured `json:"resources"`
}

// ReviewAppStatus defines the observed state of ReviewApp
type ReviewAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ReviewApp is the Schema for the reviewapps API
type ReviewApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReviewAppSpec   `json:"spec,omitempty"`
	Status ReviewAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReviewAppList contains a list of ReviewApp
type ReviewAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReviewApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ReviewApp{}, &ReviewAppList{})
}
