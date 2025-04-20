/*
Copyright 2025.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MutationTarget defines what parts of a pod can be changed
type MutationTarget string

const (
	// MutationTargetAnnotations indicates that pod annotations will be mutated
	MutationTargetAnnotations MutationTarget = "Annotations"
	// MutationTargetLabels indicates that pod labels will be mutated
	MutationTargetLabels MutationTarget = "Labels"
	// MutationTargetResources indicates that pod resource requests/limits will be mutated
	MutationTargetResources MutationTarget = "Resources"
)

// StatusType defines what type of status to track
type StatusType string

const (
	// StatusTypePodReadiness tracks pods' ready status
	StatusTypePodReadiness StatusType = "PodReadiness"
)

// AutoPodManagerSpec defines the desired state of AutoPodManager.
type AutoPodManagerSpec struct {
	// PodSelector defines which pods this manager should target
	// +required
	PodSelector *metav1.LabelSelector `json:"podSelector"`

	// Namespace specifies the namespace where pods should be selected
	// If not specified, uses the same namespace as the AutoPodManager object
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// MutationTargets specifies which parts of the pod objects can be changed
	// +required
	// +kubebuilder:validation:MinItems=1
	MutationTargets []MutationTarget `json:"mutationTargets"`

	// StatusTracking defines what kind of status to track
	// +required
	StatusTracking StatusType `json:"statusTracking"`

	// Disabled allows temporarily disabling this AutoPodManager
	// +optional
	Disabled bool `json:"disabled,omitempty"`
}

// PodStatus represents the status information for a managed pod
type PodStatus struct {
	// Name is the name of the pod
	Name string `json:"name"`

	// Ready indicates if the pod is in ready state
	Ready bool `json:"ready"`

	// LastUpdateTime shows when this status was last updated
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
}

// AutoPodManagerStatus defines the observed state of AutoPodManager.
type AutoPodManagerStatus struct {
	// ManagedPods contains status information about the pods being managed
	// +optional
	ManagedPods []PodStatus `json:"managedPods,omitempty"`

	// ReadyPodCount shows the number of pods in ready state
	// +optional
	ReadyPodCount int `json:"readyPodCount,omitempty"`

	// TotalPodCount shows the total number of pods being managed
	// +optional
	TotalPodCount int `json:"totalPodCount,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AutoPodManager is the Schema for the autopodmanagers API.
type AutoPodManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutoPodManagerSpec   `json:"spec,omitempty"`
	Status AutoPodManagerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AutoPodManagerList contains a list of AutoPodManager.
type AutoPodManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AutoPodManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AutoPodManager{}, &AutoPodManagerList{})
}

// GetPodLabelSelector converts the LabelSelector to a labels.Selector
func (a *AutoPodManager) GetPodLabelSelector() (labels.Selector, error) {
	return metav1.LabelSelectorAsSelector(a.Spec.PodSelector)
}

// ShouldManagePodLabels returns true if the AutoPodManager is configured to manage pod labels
func (a *AutoPodManager) ShouldManagePodLabels() bool {
	return containsMutationTarget(a.Spec.MutationTargets, MutationTargetLabels)
}

// ShouldManagePodAnnotations returns true if the AutoPodManager is configured to manage pod annotations
func (a *AutoPodManager) ShouldManagePodAnnotations() bool {
	return containsMutationTarget(a.Spec.MutationTargets, MutationTargetAnnotations)
}

// ShouldManagePodResources returns true if the AutoPodManager is configured to manage pod resources
func (a *AutoPodManager) ShouldManagePodResources() bool {
	return containsMutationTarget(a.Spec.MutationTargets, MutationTargetResources)
}

// containsMutationTarget checks if a specified mutation target is in the slice
func containsMutationTarget(targets []MutationTarget, target MutationTarget) bool {
	for _, t := range targets {
		if t == target {
			return true
		}
	}
	return false
}
