/*
 * Copyright (c) 2018, www.wso2.com
 */

package resources

import (
	"github.com/wso2/vick/pkg/apis/vick"
	"github.com/wso2/vick/pkg/apis/vick/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createLabels(service *v1alpha1.Service) map[string]string {
	labels := make(map[string]string, len(service.ObjectMeta.Labels)+1)

	labels[vick.ServiceNameLabelKey] = service.Name
	// order matters
	// todo: update the code if override is not possible
	for k, v := range service.ObjectMeta.Labels {
		labels[k] = v
	}
	return labels
}

func createSelector(service *v1alpha1.Service) *metav1.LabelSelector {
	return &metav1.LabelSelector{MatchLabels: createLabels(service)}
}


func deploymentName(service *v1alpha1.Service) string {
	return service.Name + "-deployment"
}


func k8sServiceName(service *v1alpha1.Service) string {
	return service.Name + "-service"
}
