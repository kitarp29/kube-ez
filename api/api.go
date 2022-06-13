package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Pod struct {
	Name            string
	Status          string
	CreatedAt       string
	UniqueID        string
	NodeName        string
	IP              string
	ContainersCount int
	ContainersInfo  []Container
	Labels          map[string]string
}

type Container struct {
	Name            string
	Image           string
	ImagePullPolicy string
	Container       int
	Port            []v1.ContainerPort
}

type Deployment struct {
	Name      string
	Status    string
	CreatedAt string
	UniqueID  string
	Labels    map[string]string
}

type Configmap struct {
	Name string
}

type Service struct {
	Name  string
	Ports string
}

func Values(UserKubeconfig string) *kubernetes.Clientset {
	log.Print("Shared Informer app started")

	kubeconfig := os.Getenv("KUBECONFIG")
	//AgentNamespace := "default"

	if kubeconfig == "" {
		kubeconfig = UserKubeconfig
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		rest.InClusterConfig()
		fmt.Printf("erorr %s building config from env\n" + err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting inclusterconfig" + err.Error())
			log.Panic(err.Error())
		}
	} else {
		log.Print("Successfully built config")
	}
	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {

		log.Panic(err.Error())
	} else {
		log.Print("Successfully built clientset")
	}

	return clientset

}

func Pods(AgentNamespace string, ContainerDetails bool) string {
	// for Pods
	clientset := Values("")

	if AgentNamespace == "" {
		AgentNamespace = "default"
	}

	var podInfo []Pod
	var containerInfo []Container
	pods, err := clientset.CoreV1().Pods(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	} else {
		for i := 0; i < len(pods.Items); i++ {

			podInfo = append(podInfo,
				Pod{
					Name:            pods.Items[i].Name,
					Status:          string(pods.Items[i].Status.Phase),
					CreatedAt:       pods.Items[i].CreationTimestamp.String(),
					UniqueID:        string(pods.Items[i].GetUID()),
					NodeName:        string(pods.Items[i].Spec.NodeName),
					IP:              string(pods.Items[i].Status.PodIP),
					ContainersCount: len(pods.Items[i].Spec.Containers),
					Labels:          pods.Items[i].Labels,
				})
			if ContainerDetails {

				for j := 0; j < len(pods.Items[i].Spec.Containers); j++ {

					containerInfo = append(containerInfo,
						Container{
							Name:            pods.Items[i].Spec.Containers[j].Name,
							Container:       j,
							Image:           pods.Items[i].Spec.Containers[j].Image,
							ImagePullPolicy: string(pods.Items[i].Spec.Containers[j].ImagePullPolicy),
							Port:            pods.Items[i].Spec.Containers[j].Ports,
						})
				}
			}
			podInfo[i].ContainersInfo = containerInfo
		}

		pods_json, err := json.Marshal(podInfo)
		if err != nil {
			log.Fatal(err)
		}

		return string(pods_json)
	}
	return "Error"
}

func Deployments(AgentNamespace string) string {
	clientset := Values("")
	if AgentNamespace == "" {
		AgentNamespace = "default"
	}

	//fmt.Printf("DEPLOYMENTS \n")
	var deploymentInfo []Deployment
	deployments, err := clientset.AppsV1().Deployments(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	} else {

		for i := 0; i < len(deployments.Items); i++ {
			//fmt.Println((deployments.Items[i].Status.Conditions))

			deploymentInfo = append(deploymentInfo,
				Deployment{
					Name:      deployments.Items[i].Name,
					Status:    string(deployments.Items[i].Status.Conditions[0].Type),
					CreatedAt: deployments.Items[i].CreationTimestamp.String(),
					UniqueID:  string(deployments.Items[i].UID),
					Labels:    deployments.Items[i].Labels,
				})
		}

		deployment_json, err := json.Marshal(deploymentInfo)
		if err != nil {
			log.Fatal(err)
		}

		return string(deployment_json)
	}
	return "Error"
}

func Configmaps(AgentNamespace string) string {
	clientset := Values("")

	if AgentNamespace == "" {
		AgentNamespace = "default"
	}

	var configmapsInfo []Configmap
	configmaps, err := clientset.CoreV1().ConfigMaps(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	} else {
		for i := 0; i < len(configmaps.Items); i++ {
			configmapsInfo = append(configmapsInfo, Configmap{configmaps.Items[i].Name})
		}

		configmap_json, err := json.Marshal(configmapsInfo)
		if err != nil {
			log.Fatal(err)
		}

		return string(configmap_json)
	}
	return "Error"
}

func Services(AgentNamespace string) string {
	clientset := Values("")

	if AgentNamespace == "" {
		AgentNamespace = "default"
	}
	var servicesInfo []Service

	services, err := clientset.CoreV1().Services(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	} else {
		for i := 0; i < len(services.Items); i++ {
			servicesInfo = append(servicesInfo, Service{Name: services.Items[i].Name, Ports: services.Items[i].Spec.Ports[0].TargetPort.String()})
		}
		service_json, err := json.Marshal(servicesInfo)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(string(pods_json))
		return string(service_json)
	}
	return "Error"
}

// onAdd is the function executed when the kubernetes informer notified the presence of a new kubernetes node in the cluster
// func onAdd(obj interface{}) {
// 	// Cast the obj as node
// 	node := obj.(*corev1.Node)
// 	_, ok := node.GetLabels()["litmus"]
// 	if ok {
// 		fmt.Printf("It has the label!\n ")
// 	} else {
// 		fmt.Printf("It does not have the label!\n")
// 	}
// }
