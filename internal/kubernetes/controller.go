package kubernetes

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/nevercase/lunara-k8s/configs"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8SController struct {
	kubeConfig *string
	clientSet  *kubernetes.Clientset
}

func NewK8SController(c *configs.Config) *K8SController {
	k := &K8SController{
		kubeConfig: &c.Kubernetes.Conf,
	}
	if k.kubeConfig == nil {
		if home := homedir.HomeDir(); home != "" {
			k.kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeConfig file")
		} else {
			k.kubeConfig = flag.String("kubeConfig", "", "absolute path to the kubeConfig file")
		}
	}
	config, err := clientcmd.BuildConfigFromFlags("", *k.kubeConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("config:", config)
	k.clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return k
}

func (k *K8SController) DeploymentList() (res *appsv1.DeploymentList) {
	deploymentsClient := k.clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return res
	}
	return list
}
