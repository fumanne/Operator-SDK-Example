/*
Copyright 2022.

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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NginxAppSpec defines the desired state of NginxApp
type NginxAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of NginxApp. Edit nginxapp_types.go to remove/update
	Num   *int32               `json:"num"`
	Image string               `json:"image"`
	Ports []corev1.ServicePort `json:"serviceports"`
	Env   []corev1.EnvVar      `json:"env"`
	Data  map[string]string    `json:"data"`
}

// NginxAppStatus defines the observed state of NginxApp
type NginxAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	appsv1.DeploymentStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NginxApp is the Schema for the nginxapps API
type NginxApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NginxAppSpec   `json:"spec"`
	Status NginxAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NginxAppList contains a list of NginxApp
type NginxAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NginxApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NginxApp{}, &NginxAppList{})
}
