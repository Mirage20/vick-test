/*
 * Copyright (c) 2018, www.wso2.com
 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Cell struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CellSpec   `json:"spec"`
	Status CellStatus `json:"status"`
}

type CellSpec struct {
	Replicas      *int32 `json:"replicas"`
	Image         string `json:"image"`
	ContainerPort int32 `json:"containerPort"`
	CellPort   int32 `json:"CellPort"`
}

type CellStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CellList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items []Cell `json:"items"`
}
