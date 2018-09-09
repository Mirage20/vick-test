/*
 * Copyright (c) 2018, www.wso2.com
 */

package resources

import (
	"github.com/wso2/vick/pkg/apis/vick/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateCoreService(service *v1alpha1.Service) *corev1.Service {
	labels := map[string]string{
		"app":        "nginx",
		"controller": service.Name,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(service, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "Service",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       "http",
				Protocol:   corev1.ProtocolTCP,
				Port:       service.Spec.ServicePort,
				TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: service.Spec.ContainerPort},
			}},
			Selector: map[string]string{
				"app": "nginx",
			},
		},
	}
}
