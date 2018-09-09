/*
 * Copyright (c) 2018, www.wso2.com
 */

package resources

import (
	"github.com/wso2/vick/pkg/apis/vick/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func CreateAppDeployment(service *v1alpha1.Service) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "nginx",
		"controller": service.Name,
	}
	podTemplateAnnotations := map[string]string{}
	podTemplateAnnotations["sidecar.istio.io/inject"] = "true"
	//https://github.com/istio/istio/blob/master/install/kubernetes/helm/istio/templates/sidecar-injector-configmap.yaml
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(service, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "Service",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: service.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: podTemplateAnnotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  service.Name,
							Image: service.Spec.Image,
							Ports:[]corev1.ContainerPort{{
								ContainerPort: service.Spec.ContainerPort,
							}},
						},
					},
				},
			},
		},
	}
}
