package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//var Kconfig chan *kubernetes.Clientset
var Kconfig *kubernetes.Clientset

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

type Event struct {
	Name       string
	Type       string
	ObjectName string
	CreatedAt  string
	UniqueID   string
}

func Main() {
	log.Print("Shared Informer app started")

	kubeconfig := os.Getenv("KUBECONFIG")

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
	//Kconfig <- clientset

	Kconfig = clientset
}

func Pods(AgentNamespace string, ContainerDetails bool) string {
	// for Pods
	clientset := Kconfig

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

func PodLogs(AgentNamespace string, PodName string) string {
	clientset := Kconfig
	req := clientset.CoreV1().Pods(AgentNamespace).GetLogs(PodName, &(v1.PodLogOptions{}))
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()

	return str
}

func Deployments(AgentNamespace string) string {
	clientset := Kconfig
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
	clientset := Kconfig

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
	clientset := Kconfig

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

func Events() string {
	clientset := Kconfig

	var eventsInfo []Event
	AgentNamespace := "default"
	events, err := clientset.CoreV1().Events(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	} else {
		for i := 0; i < len(events.Items); i++ {
			eventsInfo = append(eventsInfo,
				Event{
					Name:       events.Items[i].Name,
					ObjectName: (events.Items[i].InvolvedObject.Name),
					CreatedAt:  events.Items[i].LastTimestamp.String(),
					UniqueID:   string(events.Items[i].UID),
					Type:       events.Items[i].Type,
				})
		}
		event_json, err := json.Marshal(eventsInfo)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(string(pods_json))
		return string(event_json)
	}
	return "Error"
}
