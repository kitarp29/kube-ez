package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// setting a Global variable for the clientset so that I can resuse throughout the code
var Kconfig *kubernetes.Clientset

// These are all the Structs that are used in the API later in this code
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

type Secret struct {
	Name      string
	SecretMap map[string]string
	Type      string
	CreatedAt string
	UniqueID  string
}

type Replicationcontroller struct {
	Name      string
	CreatedAt string
	UniqueID  string
	Labels    map[string]string
}

type Daemonset struct {
	Name      string
	CreatedAt string
	UniqueID  string
	Labels    map[string]string
}

type Namespace struct {
	Name      string
	CreatedAt string
	UniqueID  string
}

type Event struct {
	Name       string
	Type       string
	ObjectName string
	CreatedAt  string
	UniqueID   string
}

//This function is used to interact with the Kubernetes Cluster to get the clienset
// It has two options:
// 1. Current it's setup to be used inside a cluster
// 2. We can configure it to be used outside the cluster

func Main() {
	logrus.Info("Shared Informer app started")

	// This checks if you have a Kubernetes config file in your home directory. If not it will try to create in in-cluster config and use that.
	var config *rest.Config
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	if err != nil {
		// If the Kubeconfig file is not available, use the in-cluster config
		logrus.Info("Using in-cluster configuration. Since couldn't find a kubeconfig file.")
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Error loading in-cluster configuration: %s\n", err)
			// BUS YHI TKK THA JO THA!!
			// So, at this point we tried to connect with local config file. Also tried to connect to one inside a cluster.
			logrus.Error(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Error(err.Error())
	}
	// comment till here

	Kconfig = clientset
}

// This function is used to get the list of all the pods in the cluster with container details
func Pods(AgentNamespace string, ContainerDetails bool, log *logrus.Entry) string {
	// for Pods
	clientset := Kconfig

	if AgentNamespace == "" {
		log.Info("Namespace is empty")
		log.Info("Namespace = default")
		AgentNamespace = "default"
	}

	var podInfo []Pod
	var containerInfo []Container
	pods, err := clientset.CoreV1().Pods(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered in Pods(): ", r)
			}
		}()
		log.Panic("Unable to find pods. Error: " + err.Error())
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
			log.Error(err.Error())
		}

		return string(pods_json)
	}
	log.Error("Error in getting pods")
	return "Error"
}

// This function is used to get the list of all the logs in a pod.
func PodLogs(AgentNamespace string, PodName string, log *logrus.Entry) string {
	clientset := Kconfig
	req := clientset.CoreV1().Pods(AgentNamespace).GetLogs(PodName, &(v1.PodLogOptions{}))
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		log.Error(err.Error())
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		log.Error(err.Error())
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()

	return str
}

// This function is used to get the list of all the deployments in the cluster
func Deployments(AgentNamespace string, log *logrus.Entry) string {
	clientset := Kconfig
	if AgentNamespace == "" {
		log.Info("Namespace is empty")
		log.Info("Namespace = default")
		AgentNamespace = "default"
	}

	//fmt.Printf("DEPLOYMENTS \n")
	var deploymentInfo []Deployment
	deployments, err := clientset.AppsV1().Deployments(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered in Deployments(): ", r)
			}
		}()
		log.Panic("Unable to find Deployments. Error: " + err.Error())
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
			log.Error(err.Error())
			log.Fatal(err)
		}

		return string(deployment_json)
	}
	log.Error("Error in getting deployments")
	return "Error"
}

// This function is used to get the list of all the Configmaps in the cluster
func Configmaps(AgentNamespace string, log *logrus.Entry) string {
	clientset := Kconfig

	if AgentNamespace == "" {
		log.Info("Namespace is empty")
		log.Info("Namespace = default")
		AgentNamespace = "default"
	}

	var configmapsInfo []Configmap
	configmaps, err := clientset.CoreV1().ConfigMaps(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered in Configmaps(): ", r)
			}
		}()
		log.Panic("Unable to find Configmaps. Error: " + err.Error())
	} else {
		for i := 0; i < len(configmaps.Items); i++ {
			configmapsInfo = append(configmapsInfo, Configmap{configmaps.Items[i].Name})
		}

		configmap_json, err := json.Marshal(configmapsInfo)
		if err != nil {
			log.Print(err.Error())
			log.Fatal(err)
		}

		return string(configmap_json)
	}
	log.Error("Error in getting configmaps")
	return "Error"
}

// This function is used to get the list of all the Services in the cluster
func Services(AgentNamespace string, log *logrus.Entry) string {
	clientset := Kconfig

	if AgentNamespace == "" {
		log.Info("Namespace is empty")
		log.Info("Namespace = default")
		AgentNamespace = "default"
	}
	var servicesInfo []Service

	services, err := clientset.CoreV1().Services(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered in Services()", r)
			}
		}()
		log.Panic("Unable to find Services. Error: " + err.Error())
	} else {
		for i := 0; i < len(services.Items); i++ {
			servicesInfo = append(servicesInfo, Service{Name: services.Items[i].Name, Ports: services.Items[i].Spec.Ports[0].TargetPort.String()})
		}
		service_json, err := json.Marshal(servicesInfo)
		if err != nil {
			log.Error(err.Error())
			log.Fatal(err)
		}
		//fmt.Println(string(pods_json))
		return string(service_json)
	}
	log.Error("Error in getting services")
	return "Error"
}

// This function is used to get the list of all the events in the cluster
func Events(AgentNamespace string, log *logrus.Entry) string {
	clientset := Kconfig

	var eventsInfo []Event
	if AgentNamespace == "" {
		log.Info("Namespace is empty")
		log.Info("Namespace = default")
		AgentNamespace = "default"
	}
	events, err := clientset.CoreV1().Events(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered in Events()", r)
			}
		}()
		log.Panic("Unable to find events. Error: " + err.Error())
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
			defer func() {
				if r := recover(); r != nil {
					log.Error("Recovered in Events()", r)
				}
			}()
			log.Panic("Issue in marshaling JSON in Events. Error: " + err.Error())
		}
		//fmt.Println(string(pods_json))
		return string(event_json)
	}
	log.Error("Error in getting events")
	return "Error"
}

// This function is used to get the list of all the secrets in the cluster
func Secrets(AgentNamespace string, log *logrus.Entry) string {
	clientset := Kconfig
	if AgentNamespace == "" {
		log.Info("Namespace is empty")
		log.Info("Namespace = default")
		AgentNamespace = "default"
	}
	secrets, err := clientset.CoreV1().Secrets(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered in Secrets()", r)
			}
		}()
		log.Panic("Unable to find secrets. Error: " + err.Error())
	} else {
		var secretInfo []Secret
		for i := 0; i < len(secrets.Items); i++ {
			secretInfo = append(secretInfo,
				Secret{
					Name:      secrets.Items[i].Name,
					Type:      string(secrets.Items[i].Type),
					CreatedAt: secrets.Items[i].CreationTimestamp.String(),
					UniqueID:  string(secrets.Items[i].UID),
				})
			tmp := make(map[string]string)
			for key, value := range secrets.Items[i].Data {
				tmp[key] = string(value)
			}
			secretInfo[i].SecretMap = tmp
		}
		secret_json, err := json.Marshal(secretInfo)
		if err != nil {
			defer func() {
				if r := recover(); r != nil {
					log.Error("Recovered in Secrets()", r)
				}
			}()
			log.Panic("Error in Marshalling JSON in Secrets. Error: " + err.Error())
		}
		//fmt.Println(string(secret_json))
		return string(secret_json)
	}
	log.Error("Error in getting secrets")
	return "Error"
}

// This function is used to get the list of all the ReplicaController in the cluster
func ReplicationController(AgentNamespace string, log *logrus.Entry) string {
	clientset := Kconfig
	if AgentNamespace == "" {
		log.Info("Namespace is empty")
		log.Info("Namespace = default")
		AgentNamespace = "default"
	}
	replicationcontrollers, err := clientset.CoreV1().ReplicationControllers(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered in ReplicaController()", r)
			}
		}()
		log.Panic("Unable to find ReplicaControllers. Error: " + err.Error())
	} else {
		var replicationcontrollerInfo []Replicationcontroller
		for i := 0; i < len(replicationcontrollers.Items); i++ {
			replicationcontrollerInfo = append(replicationcontrollerInfo,
				Replicationcontroller{
					Name:      replicationcontrollers.Items[i].Name,
					CreatedAt: replicationcontrollers.Items[i].CreationTimestamp.String(),
					UniqueID:  string(replicationcontrollers.Items[i].UID),
					Labels:    (replicationcontrollers.Items[i].Labels),
				})
		}
		replicationcontroller_json, err := json.Marshal(replicationcontrollerInfo)
		if err != nil {
			defer func() {
				if r := recover(); r != nil {
					log.Error("Recovered in ReplicaController()", r)
				}
			}()
			log.Panic("Error in Marshalling JSON in ReplicaController. Error: " + err.Error())
		}
		//fmt.Println(string(replicationcontroller_json))
		return string(replicationcontroller_json)
	}
	log.Error("Error in getting replicationcontrollers")
	return "Error"
}

// This function is used to get the list of all the Daemonsets in the cluster
func DaemonSet(AgentNamespace string, log *logrus.Entry) string {
	clientset := Kconfig
	if AgentNamespace == "" {
		log.Info("Namespace is empty")
		log.Info("Namespace = default")
		AgentNamespace = "default"
	}
	daemonsets, err := clientset.ExtensionsV1beta1().DaemonSets(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
	} else {
		var daemonsetInfo []Daemonset
		for i := 0; i < len(daemonsets.Items); i++ {
			daemonsetInfo = append(daemonsetInfo,
				Daemonset{
					Name:      daemonsets.Items[i].Name,
					CreatedAt: daemonsets.Items[i].CreationTimestamp.String(),
					UniqueID:  string(daemonsets.Items[i].UID),
					Labels:    (daemonsets.Items[i].Labels),
				})
		}
		daemonset_json, err := json.Marshal(daemonsetInfo)
		if err != nil {
			log.Error(err.Error())
		}
		//fmt.Println(string(daemonset_json))
		return string(daemonset_json)
	}
	log.Error("Error in getting daemonsets")
	return "Error"
}

// This function is used to get the list of all the Namespaces in the cluster
func NameSpace(log *logrus.Entry) string {
	clientset := Kconfig
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered in NameSpace()", r)
			}
		}()
		log.Panic("Unable to find namespaces. Error: " + err.Error())
	} else {
		var namespaceInfo []Namespace
		for i := 0; i < len(namespaces.Items); i++ {
			namespaceInfo = append(namespaceInfo,
				Namespace{
					Name:      namespaces.Items[i].Name,
					CreatedAt: namespaces.Items[i].CreationTimestamp.String(),
					UniqueID:  string(namespaces.Items[i].UID),
				})
		}
		namespace_json, err := json.Marshal(namespaceInfo)
		if err != nil {
			log.Error(err.Error())
			log.Fatal(err)
		}
		return string(namespace_json)
	}
	log.Error("Error in getting namespaces")
	return "Error"
}

// This function creates Namespace in the cluster
func CreateNamespace(namespace string, log *logrus.Entry) string {
	log.Info("Namespace=" + namespace)
	clientset := Kconfig
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				"name": namespace,
			},
		},
	}
	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	log.Info("Namespace" + namespace + "successfully")
	return "Namespace: " + namespace + " Created!"
}

// This function deletes Namespace in the cluster
func DeleteNamespace(namespace string, log *logrus.Entry) string {
	clientset := Kconfig
	err := clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	log.Info("Namespace: " + namespace + " Deleted!")
	return "Namespace: " + namespace + " Deleted!"
}

// This function Deletes the Deployments
func DeleteDeployment(namespace string, deployment string, log *logrus.Entry) string {
	clientset := Kconfig
	err := clientset.AppsV1().Deployments(namespace).Delete(context.Background(), deployment, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	log.Info("Deployment: " + deployment + " Deleted!")
	return "Deployment: " + deployment + " Deleted!"
}

// This function Deletes the services
func DeleteService(namespace string, service string, log *logrus.Entry) string {
	clientset := Kconfig
	err := clientset.CoreV1().Services(namespace).Delete(context.Background(), service, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	log.Info("Service: " + service + " Deleted!")
	return "Service: " + service + " Deleted!"
}

// This function Deletes the ConfigMap
func DeleteConfigMap(namespace string, configmap string, log *logrus.Entry) string {
	clientset := Kconfig
	err := clientset.CoreV1().ConfigMaps(namespace).Delete(context.Background(), configmap, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	log.Info("ConfigMap: " + configmap + " Deleted!")
	return "ConfigMap: " + configmap + " Deleted!"
}

// This function Deletes the Secrets
func DeleteSecret(namespace string, secret string, log *logrus.Entry) string {
	clientset := Kconfig
	err := clientset.CoreV1().Secrets(namespace).Delete(context.Background(), secret, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	log.Info("Secret: " + secret + " Deleted!")
	return "Secret: " + secret + " Deleted!"
}

// This function Deletes the ReplicationController
func DeleteReplicationController(namespace string, replicationcontroller string, log *logrus.Entry) string {
	clientset := Kconfig
	err := clientset.CoreV1().ReplicationControllers(namespace).Delete(context.Background(), replicationcontroller, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	log.Info("ReplicationController: " + replicationcontroller + " Deleted!")
	return "ReplicationController: " + replicationcontroller + " Deleted!"
}

// This function Deletes the DaemonSet
func DeleteDaemonSet(namespace string, daemonset string, log *logrus.Entry) string {
	clientset := Kconfig
	err := clientset.ExtensionsV1beta1().DaemonSets(namespace).Delete(context.Background(), daemonset, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	log.Info("DaemonSet: " + daemonset + " Deleted!")
	return "DaemonSet: " + daemonset + " Deleted!"
}

// This function Deletes the Pod
func DeletePod(namespace string, pod string) string {
	clientset := Kconfig
	err := clientset.CoreV1().Pods(namespace).Delete(context.Background(), pod, metav1.DeleteOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	return "Pod: " + pod + " Deleted!"
}

// This function Deletes the Event
func DeleteEvent(namespace string, event string, log *logrus.Entry) string {
	clientset := Kconfig
	err := clientset.CoreV1().Events(namespace).Delete(context.Background(), event, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	log.Info("Event: " + event + " Deleted!")
	return "Event: " + event + " Deleted!"
}

// This function Deletes EVERYTHING in the namespace. My lil nuke!! MUWAHAHAHA
func DeleteAll(namespace string, log *logrus.Entry) string {
	clientset := Kconfig
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	for i := 0; i < len(deployments.Items); i++ {
		err := clientset.AppsV1().Deployments(namespace).Delete(context.Background(), deployments.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Error(err.Error())
			return err.Error()
		}
	}
	services, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	for i := 0; i < len(services.Items); i++ {
		err := clientset.CoreV1().Services(namespace).Delete(context.Background(), services.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Error(err.Error())
			return err.Error()
		}
	}
	configmaps, err := clientset.CoreV1().ConfigMaps(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	for i := 0; i < len(configmaps.Items); i++ {
		err := clientset.CoreV1().ConfigMaps(namespace).Delete(context.Background(), configmaps.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Error(err.Error())
			return err.Error()
		}
	}
	secrets, err := clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	for i := 0; i < len(secrets.Items); i++ {
		err := clientset.CoreV1().Secrets(namespace).Delete(context.Background(), secrets.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Error(err.Error())
			return err.Error()
		}
	}
	replicationcontrollers, err := clientset.CoreV1().ReplicationControllers(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	for i := 0; i < len(replicationcontrollers.Items); i++ {
		err := clientset.CoreV1().ReplicationControllers(namespace).Delete(context.Background(), replicationcontrollers.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Error(err.Error())
			return err.Error()
		}
	}
	daemonsets, err := clientset.ExtensionsV1beta1().DaemonSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	for i := 0; i < len(daemonsets.Items); i++ {
		err := clientset.ExtensionsV1beta1().DaemonSets(namespace).Delete(context.Background(), daemonsets.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Error(err.Error())
			return err.Error()
		}
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	for i := 0; i < len(pods.Items); i++ {
		err := clientset.CoreV1().Pods(namespace).Delete(context.Background(), pods.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Error(err.Error())
			return err.Error()
		}
	}
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	for i := 0; i < len(events.Items); i++ {
		err := clientset.CoreV1().Events(namespace).Delete(context.Background(), events.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Error(err.Error())
			return err.Error()
		}
	}
	log.Info("Everything in " + namespace + " Deleted!")
	return "All Deleted!"
}
