package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

func Main() {
	log.Print("Shared Informer app started")
	// This will be used in case you have to run the code outside the cluster
	// You will have to export the KUBECONFIG variable to point to the config file in the terminal
	// kubeconfig := os.Getenv("KUBECONFIG")
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// if err != nil {
	// 	rest.InClusterConfig()
	// 	fmt.Printf("erorr %s building config from env\n" + err.Error())
	// 	config, err = rest.InClusterConfig()
	// 	if err != nil {
	// 		fmt.Printf("error %s, getting inclusterconfig" + err.Error())
	// 		log.Print(err.Error())
	// 		log.Panic(err.Error())
	// 	}
	// } else {
	// 	log.Print("Successfully built config")
	// }
	// // Create the clientset
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	log.Print(err.Error())
	// 	log.Panic(err.Error())
	// } else {
	// 	log.Print("Successfully built clientset")
	// }

	// uncomment this if you want to learn this file inside the cluster
	// the file above this has to be run inside the cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

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
		log.Print(err.Error())
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
			log.Print(err.Error())
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
		log.Print(err.Error())
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
		log.Print(err.Error())
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
			log.Print(err.Error())
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
		log.Print(err.Error())
		log.Panic(err.Error())
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
		log.Print(err.Error())
		log.Panic(err.Error())
	} else {
		for i := 0; i < len(services.Items); i++ {
			servicesInfo = append(servicesInfo, Service{Name: services.Items[i].Name, Ports: services.Items[i].Spec.Ports[0].TargetPort.String()})
		}
		service_json, err := json.Marshal(servicesInfo)
		if err != nil {
			log.Print(err.Error())
			log.Fatal(err)
		}
		//fmt.Println(string(pods_json))
		return string(service_json)
	}
	return "Error"
}

func Events(AgentNamespace string) string {
	clientset := Kconfig

	var eventsInfo []Event
	if AgentNamespace == "" {
		AgentNamespace = "default"
	}
	events, err := clientset.CoreV1().Events(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
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
			log.Print(err.Error())
			log.Fatal(err)
		}
		//fmt.Println(string(pods_json))
		return string(event_json)
	}
	return "Error"
}
func Secrets(AgentNamespace string) string {
	clientset := Kconfig
	if AgentNamespace == "" {
		AgentNamespace = "default"
	}
	secrets, err := clientset.CoreV1().Secrets(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		log.Panic(err.Error())
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
			log.Print(err.Error())
			log.Fatal(err)
		}
		//fmt.Println(string(secret_json))
		return string(secret_json)
	}
	return "Error"
}

func ReplicationController(AgentNamespace string) string {
	clientset := Kconfig
	if AgentNamespace == "" {
		AgentNamespace = "default"
	}
	replicationcontrollers, err := clientset.CoreV1().ReplicationControllers(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		log.Panic(err.Error())
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
			log.Print(err.Error())
			log.Fatal(err)
		}
		//fmt.Println(string(replicationcontroller_json))
		return string(replicationcontroller_json)
	}
	return "Error"
}

func DaemonSet(AgentNamespace string) string {
	clientset := Kconfig
	if AgentNamespace == "" {
		AgentNamespace = "default"
	}
	daemonsets, err := clientset.ExtensionsV1beta1().DaemonSets(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		log.Panic(err.Error())
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
			log.Print(err.Error())
			log.Panic(err)
		}
		//fmt.Println(string(daemonset_json))
		return string(daemonset_json)
	}
	return "Error"
}

func NameSpace() string {
	clientset := Kconfig
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		log.Panic(err.Error())
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
			log.Print(err.Error())
			log.Fatal(err)
		}
		return string(namespace_json)
	}
	return "Error"
}

func CreateNamespace(namespace string) string {
	fmt.Println(namespace)
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
		log.Print(err.Error())
		return err.Error()
	}
	return "Namespace: " + namespace + " Created!"
}

func DeleteNamespace(namespace string) string {
	clientset := Kconfig
	err := clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	return "Namespace: " + namespace + " Deleted!"
}

func DeleteDeployment(namespace string, deployment string) string {
	clientset := Kconfig
	err := clientset.AppsV1().Deployments(namespace).Delete(context.Background(), deployment, metav1.DeleteOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	return "Deployment: " + deployment + " Deleted!"
}

func DeleteService(namespace string, service string) string {
	clientset := Kconfig
	err := clientset.CoreV1().Services(namespace).Delete(context.Background(), service, metav1.DeleteOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	return "Service: " + service + " Deleted!"
}

func DeleteConfigMap(namespace string, configmap string) string {
	clientset := Kconfig
	err := clientset.CoreV1().ConfigMaps(namespace).Delete(context.Background(), configmap, metav1.DeleteOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	return "ConfigMap: " + configmap + " Deleted!"
}

func DeleteSecret(namespace string, secret string) string {
	clientset := Kconfig
	err := clientset.CoreV1().Secrets(namespace).Delete(context.Background(), secret, metav1.DeleteOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	return "Secret: " + secret + " Deleted!"
}

func DeleteReplicationController(namespace string, replicationcontroller string) string {
	clientset := Kconfig
	err := clientset.CoreV1().ReplicationControllers(namespace).Delete(context.Background(), replicationcontroller, metav1.DeleteOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	return "ReplicationController: " + replicationcontroller + " Deleted!"
}

func DeleteDaemonSet(namespace string, daemonset string) string {
	clientset := Kconfig
	err := clientset.ExtensionsV1beta1().DaemonSets(namespace).Delete(context.Background(), daemonset, metav1.DeleteOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	return "DaemonSet: " + daemonset + " Deleted!"
}

func DeletePod(namespace string, pod string) string {
	clientset := Kconfig
	err := clientset.CoreV1().Pods(namespace).Delete(context.Background(), pod, metav1.DeleteOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	return "Pod: " + pod + " Deleted!"
}

func DeleteEvent(namespace string, event string) string {
	clientset := Kconfig
	err := clientset.CoreV1().Events(namespace).Delete(context.Background(), event, metav1.DeleteOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	return "Event: " + event + " Deleted!"
}

func DeleteAll(namespace string) string {
	clientset := Kconfig
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	for i := 0; i < len(deployments.Items); i++ {
		err := clientset.AppsV1().Deployments(namespace).Delete(context.Background(), deployments.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
	}
	services, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	for i := 0; i < len(services.Items); i++ {
		err := clientset.CoreV1().Services(namespace).Delete(context.Background(), services.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
	}
	configmaps, err := clientset.CoreV1().ConfigMaps(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	for i := 0; i < len(configmaps.Items); i++ {
		err := clientset.CoreV1().ConfigMaps(namespace).Delete(context.Background(), configmaps.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
	}
	secrets, err := clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	for i := 0; i < len(secrets.Items); i++ {
		err := clientset.CoreV1().Secrets(namespace).Delete(context.Background(), secrets.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Print(err.Error())
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
			log.Print(err.Error())
			return err.Error()
		}
	}
	daemonsets, err := clientset.ExtensionsV1beta1().DaemonSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	for i := 0; i < len(daemonsets.Items); i++ {
		err := clientset.ExtensionsV1beta1().DaemonSets(namespace).Delete(context.Background(), daemonsets.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	for i := 0; i < len(pods.Items); i++ {
		err := clientset.CoreV1().Pods(namespace).Delete(context.Background(), pods.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
	}
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	}
	for i := 0; i < len(events.Items); i++ {
		err := clientset.CoreV1().Events(namespace).Delete(context.Background(), events.Items[i].Name, metav1.DeleteOptions{})
		if err != nil {
			log.Print(err.Error())
			return err.Error()
		}
	}
	return "All Deleted!"
}
