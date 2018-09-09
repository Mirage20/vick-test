/*
 * Copyright (c) 2018, www.wso2.com
 */

package service

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/wso2/vick/pkg/controller"
	"github.com/wso2/vick/pkg/controller/service/resources"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	appsv1informers "k8s.io/client-go/informers/apps/v1"
	corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	//corev1informers "k8s.io/client-go/informers/core/v1"
	vickinformers "github.com/wso2/vick/pkg/client/informers/externalversions/vick/v1alpha1"
	listers "github.com/wso2/vick/pkg/client/listers/vick/v1alpha1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

type serviceHandler struct {
	serviceLister    listers.ServiceLister
	deploymentLister appsv1listers.DeploymentLister
	k8sServiceLister corev1listers.ServiceLister
	kubeclientset    kubernetes.Interface
}

func NewController(kubeClient kubernetes.Interface, k8sServiceInformer corev1informers.ServiceInformer, serviceInformer vickinformers.ServiceInformer, deploymentInformer appsv1informers.DeploymentInformer) *controller.Controller {
	h := &serviceHandler{
		kubeclientset:    kubeClient,
		serviceLister:    serviceInformer.Lister(),
		k8sServiceLister: k8sServiceInformer.Lister(),
		deploymentLister: deploymentInformer.Lister(),
	}
	c := controller.New(h, "Service")

	glog.Info("Setting up event handlers")
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.Enqueue,
		UpdateFunc: func(old, new interface{}) {
			glog.Infof("Old %+v\nnew %+v", old, new)
			c.Enqueue(new)
		},
		DeleteFunc: c.Enqueue,
	})
	return c
}

func (h *serviceHandler) Handle(key string) error {
	glog.Infof("Handle called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		glog.Errorf("invalid resource key: %s", key)
		return nil
	}
	service, err := h.serviceLister.Services(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("service '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	glog.Infof("Found service %+v", service)
	// Get the deployment with the name specified in Foo.spec
	deployment, err := h.deploymentLister.Deployments(service.Namespace).Get(service.Name)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = h.kubeclientset.AppsV1().Deployments(service.Namespace).Create(resources.CreateAppDeployment(service))
	}

	k8sService, err := h.k8sServiceLister.Services(service.Namespace).Get(service.Name)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		k8sService, err = h.kubeclientset.CoreV1().Services(service.Namespace).Create(resources.CreateCoreService(service))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	glog.Infof("Deployment created %+v", deployment)
	glog.Infof("Service created %+v", k8sService)
	// If the Deployment is not controlled by this Foo resource, we should log
	// a warning to the event recorder and ret
	//if !metav1.IsControlledBy(deployment, foo) {
	//	msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
	//	c.recorder.Event(foo, corev1.EventTypeWarning, ErrResourceExists, msg)
	//	return fmt.Errorf(msg)
	//}

	// If this number of the replicas on the Foo resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	//if foo.Spec.Replicas != nil && *foo.Spec.Replicas != *deployment.Spec.Replicas {
	//	glog.V(4).Infof("Foo %s replicas: %d, deployment replicas: %d", name, *foo.Spec.Replicas, *deployment.Spec.Replicas)
	//	deployment, err = c.kubeclientset.AppsV1().Deployments(foo.Namespace).Update(newDeployment(foo))
	//}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	//if err != nil {
	//	return err
	//}

	return nil
}