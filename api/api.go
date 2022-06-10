package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Pod struct {
	Name   string
	Status string
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

func Pods(AgentNamespace string) string {
	// for Pods
	clientset := Values("")
	fmt.Printf("PODS \n")
	var podInfo []Pod
	pods, err := clientset.CoreV1().Pods(AgentNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	} else {
		for i := 0; i < len(pods.Items); i++ {
			podInfo = append(podInfo, Pod{pods.Items[i].Name, string(pods.Items[i].Status.Phase)})
		}
		fmt.Printf("%v\n", podInfo)
		pods_json, err := json.Marshal(podInfo)

		if err != nil {

			log.Fatal(err)
		}

		fmt.Println(string(pods_json))

		return string(pods_json)
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
