/*
 * Copyright (c) 2018, www.wso2.com
 */

package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	vickclientset "github.com/wso2/vick/pkg/client/clientset/versioned"
	vickinformers "github.com/wso2/vick/pkg/client/informers/externalversions"
	"github.com/wso2/vick/pkg/controller/cell"
	"github.com/wso2/vick/pkg/controller/service"
	"github.com/wso2/vick/pkg/signals"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	threadsPerController = 2
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	vickClient, err := vickclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building vick clientset: %s", err.Error())
	}

	// Create informers
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	vickInformerFactory := vickinformers.NewSharedInformerFactory(vickClient, time.Second*30)

	deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
	k8sServiceInformer := kubeInformerFactory.Core().V1().Services()
	serviceInformer := vickInformerFactory.Vick().V1alpha1().Services()
	cellInformer := vickInformerFactory.Vick().V1alpha1().Cells()

	// Create crd controllers
	cellController := cell.NewController(kubeClient, cellInformer)
	serviceController := service.NewController(kubeClient,vickClient, k8sServiceInformer, cellInformer, serviceInformer, deploymentInformer)

	// Start informers
	go kubeInformerFactory.Start(stopCh)
	go vickInformerFactory.Start(stopCh)

	// Wait for cache sync
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh,
		deploymentInformer.Informer().HasSynced,
		k8sServiceInformer.Informer().HasSynced,
		cellInformer.Informer().HasSynced,
		serviceInformer.Informer().HasSynced); !ok {
		glog.Fatal("failed to wait for caches to sync")
	}

	//Start controllers
	go cellController.Run(threadsPerController, stopCh)
	go serviceController.Run(threadsPerController, stopCh)
	<-stopCh
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
