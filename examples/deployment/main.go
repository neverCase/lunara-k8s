package main

import (
	"flag"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
)

func main() {
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	log.Println("host:", host, " port:", port)

	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeConfig file")
	} else {
		kubeConfig = flag.String("kubeConfig", "", "absolute path to the kubeConfig file")
	}
	flag.Parse()
	log.Println("kubeConfig:", *kubeConfig)
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err)
	}
	log.Println("config:", config)
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}
