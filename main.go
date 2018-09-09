package main

import (
	"flag"

	"time"

	"github.com/golang/glog"
	clientset "github.com/wso2/vick/pkg/client/clientset/versioned"
	vickinformers "github.com/wso2/vick/pkg/client/informers/externalversions"
	"github.com/wso2/vick/pkg/signals"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

var (
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", os.Getenv("HOME")+"/.kube/config", "kubeconfig path. Default $HOME/.kube/config.")
}

func main() {
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	vickClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	vickInformerFactory := vickinformers.NewSharedInformerFactory(vickClient, time.Second*30)

	controller := NewController(kubeClient, vickClient,
		vickInformerFactory.Vickcontroller().V1alpha1().Cells())

	go kubeInformerFactory.Start(stopCh)
	go vickInformerFactory.Start(stopCh)

	if err = controller.Run(1, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
