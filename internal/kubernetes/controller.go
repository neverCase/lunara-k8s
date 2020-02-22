package kubernetes

import (
	"log"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Config struct {
	Conf      string `json:"conf" yaml:"conf"`
	MasterUrl string `json:"master_url" yaml:"master_url"`
}

type K8SController struct {
	kubeConfig Config
	ClientSet  *kubernetes.Clientset
}

func NewK8SController(c Config) *K8SController {
	k := &K8SController{
		kubeConfig: c,
	}
	if k.kubeConfig.Conf == "" && k.kubeConfig.MasterUrl == "" {
		k.kubeConfig.Conf = filepath.Join(homedir.HomeDir(), ".kube", "config")
		k.kubeConfig.MasterUrl = ""
	}
	config, err := clientcmd.BuildConfigFromFlags(k.kubeConfig.MasterUrl, k.kubeConfig.Conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("config:", config)
	k.ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return k
}

func (k *K8SController) DeploymentList() (res *appsv1.DeploymentList) {
	deploymentsClient := k.ClientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return res
	}
	return list
}
