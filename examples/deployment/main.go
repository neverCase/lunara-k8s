package main

import (
	"fmt"

	"github.com/json-iterator/go"
	k8sCrd "github.com/nevercase/lunara-k8s/internal/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	//log.Println("host:", host, " port:", port)
	//
	//var kubeConfig *string
	//if home := homedir.HomeDir(); home != "" {
	//	kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeConfig file")
	//} else {
	//	kubeConfig = flag.String("kubeConfig", "", "absolute path to the kubeConfig file")
	//}
	//flag.Parse()
	//log.Println("kubeConfig:", *kubeConfig)
	//config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	//if err != nil {
	//	panic(err)
	//}
	//log.Println("config:", config)
	//clientSet, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	panic(err)
	//}
	k := k8sCrd.NewK8SController(k8sCrd.Config{})

	deploymentsClient := k.ClientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-redis",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "test-redis",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "test-redis",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "tr-container",
							Image: "39.99.229.222:37229/helix-saga/test-redis:2.0",
							Ports: []apiv1.ContainerPort{
								{
									//Name:          "http",
									//Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 81,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "GET_HOSTS_FROM",
									Value: "dns",
								},
							},
						},
					},
					ImagePullSecrets: []apiv1.LocalObjectReference{
						{
							Name: "private-harbor",
						},
					},
					RestartPolicy: apiv1.RestartPolicyAlways,
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	type DataList struct {
		Deployments appsv1.DeploymentList `json:"deployments"`
	}
	type Res struct {
		ErrorCode int    `json:"error_code"`
		Message   string `json:"message"`
		Data      string `json:"data"`
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	res := &Res{
		ErrorCode: 10000,
		Message:   "success",
	}
	data := DataList{
		Deployments: *list,
	}
	content, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	res.Data = string(content)
	response, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	fmt.Println("response-1:", string(response))

	//dynamicDemo(kubeConfig)
}

func dynamicDemo(kubeConfig *string) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err)
	}
	namespace := "default"
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "test-redis",
			},
			"spec": map[string]interface{}{
				"replicas": 1,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "demo",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "demo",
						},
					},

					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  "tr-container",
								"image": "39.99.229.222:37229/helix-saga/test-redis:1.0",
								"ports": []map[string]interface{}{
									{
										//"name":          "http",
										//"protocol":      "TCP",
										"containerPort": 81,
									},
								},
								"env": []map[string]interface{}{
									{
										"name":  "GET_HOSTS_FROM",
										"value": "dns",
									},
								},
							},
						},
						"imagePullSecrets": []map[string]interface{}{
							{
								"name": "private-harbor",
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := client.Resource(deploymentRes).Namespace(namespace).Create(deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetName())
}

func int32Ptr(i int32) *int32 { return &i }
