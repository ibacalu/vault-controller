/*
Copyright 2023.

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

// PolicySpec defines the desired state of Policy
type PolicySpec struct {
	// Policies provide a declarative way to grant or forbid access to certain paths and operations in Vault
	Rules []PolicyRule `json:"rules"`
}

// PolicyRule represents a single rule in the policy.
type PolicyRule struct {
	// Policies use path-based matching to test the set of capabilities against a request. A policy path may specify an exact path to match, or it could specify a glob pattern which instructs Vault to use a prefix match
	Path string `json:"path"`
	// Provide fine-grained control over permitted (or denied) operations. Options are: [create, read, update, patch, delete, list]
	Capabilities []string `json:"capabilities"`
}

// PolicyStatus defines the observed state of Policy
type PolicyStatus struct {
	// Was the policy applied successfully?
	Applied bool `json:"applied"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Policy is the Schema for the policies API
type Policy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PolicySpec   `json:"spec,omitempty"`
	Status PolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PolicyList contains a list of Policy
type PolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Policy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Policy{}, &PolicyList{})
}
