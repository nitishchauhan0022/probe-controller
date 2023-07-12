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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProbeSpec defines the desired state of Probe
type ProbeSpec struct {
	Selector      metav1.LabelSelector `json:"selector"`
	Probe         corev1.Probe         `json:"probe"`
	ContainerName string               `json:"containerName,omitempty"`
	SlackDetails  SlackDetails         `json:"slackDetails"`
}

// ProbeStatus defines the observed state of Probe
type ProbeStatus struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	FailureCount int    `json:"failureCount"`
}

// slackDetails defines the slack details
type SlackDetails struct {
	SlackToken   ResolveFromSecret `json:"slackToken"`
	SlackChannel string            `json:"slackChannel"`
}

// resolveFromSecret resolves the secre
type ResolveFromSecret struct {
	SecretName string `json:"secretName"`
	SecretKey  string `json:"secretKey"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Probe is the Schema for the probes API
type Probe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProbeSpec   `json:"spec,omitempty"`
	Status ProbeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProbeList contains a list of Probe
type ProbeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Probe `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Probe{}, &ProbeList{})
}
