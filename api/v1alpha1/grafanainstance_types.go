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
	appsv1 "k8s.io/api/apps/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GrafanaInstanceSpec defines the desired state of GrafanaInstance
type GrafanaInstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Image                 string                          `json:"image,omitempty"`
	Port                  int32                           `json:"port,omitempty"`
	CredentialsSecretName string                          `json:"credentialsSecretName,omitempty"`
	INIConfig             map[string]apiextensionsv1.JSON `json:"iniConfig,omitempty"`
}

// GrafanaInstanceStatus defines the observed state of GrafanaInstance
type GrafanaInstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	GrafanaUI GrafanaUIStatus `json:"grafanaUI,omitempty"`
}

type GrafanaUIStatus struct {
	AvailableReplicas string                       `json:"availableReplicas,omitempty"`
	Conditions        []appsv1.DeploymentCondition `json:"conditions,omitempty"`
	ServiceURL        string                       `json:"serviceURL,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GrafanaInstance is the Schema for the grafanainstances API
type GrafanaInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrafanaInstanceSpec   `json:"spec,omitempty"`
	Status GrafanaInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GrafanaInstanceList contains a list of GrafanaInstance
type GrafanaInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrafanaInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GrafanaInstance{}, &GrafanaInstanceList{})
}
