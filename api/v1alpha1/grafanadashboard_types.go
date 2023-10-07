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
	"time"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	GrafanaGeneralFolder = "General"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GrafanaDashboardSpec defines the desired state of GrafanaDashboard
type GrafanaDashboardSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	Json apiextensionsv1.JSON `json:"json,omitempty"`

	// +optional
	Folder string `json:"folder,omitempty"`

	// +optional
	Name string `json:"name"`

	// SyncPeriod is the time duration to wait between each sync operation.
	// The operator will check the actual state in Grafana and reconcile it with the desired state defined in the custom resource.
	// +optional
	SyncPeriod metav1.Duration `json:"syncPeriod,omitempty"`

	// Reference to the GrafanaInstance that this dashboard should be associated with
	GrafanaInstanceRef GrafanaInstanceRef `json:"grafanaInstanceRef"`
}

// Set default values
func (d *GrafanaDashboard) SetDefaults() {
	if d.Spec.SyncPeriod.Duration == 0 {
		d.Spec.SyncPeriod.Duration = 5 * time.Minute
	}

	if d.Spec.Folder == "" {
		d.Spec.Folder = GrafanaGeneralFolder
	}
}

// GrafanaInstanceRef defines the reference to a GrafanaInstance
type GrafanaInstanceRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// GrafanaDashboardStatus defines the observed state of GrafanaDashboard
type GrafanaDashboardStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	FolderUID    string `json:"folderUID,omitempty"`
	DashboardUID string `json:"dashboardUID,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GrafanaDashboard is the Schema for the grafanadashboards API
type GrafanaDashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrafanaDashboardSpec   `json:"spec,omitempty"`
	Status GrafanaDashboardStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GrafanaDashboardList contains a list of GrafanaDashboard
type GrafanaDashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrafanaDashboard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GrafanaDashboard{}, &GrafanaDashboardList{})
}
