# **K8s-API**
[![GitHub contributors](https://img.shields.io/github/contributors/kitarp29/k8s-api)](https://github.com/GDSC-DSI/Schedura/graphs/contributors) [![GitHub issues](https://img.shields.io/github/issues/HAC-2020/Aimers)](https://github.com/kitarp29/k8s-api/issues/) 
![GitHub stars](https://img.shields.io/github/stars/kitarp29/k8s-api) ![GitHub releases](https://img.shields.io/github/release/kitarp29/k8s-api)![GitHub license](https://img.shields.io/github/license/kitarp29/k8s-api)
![Docker](https://img.shields.io/docker/pulls/kitarp29/k8s-api)
![Docker](https://img.shields.io/docker/stars/kitarp29/k8s-api)
![Follow](https://img.shields.io/twitter/follow/kitarp29?label=Follow+Kitarp29&style=social)

<img src="golang.jpeg"></img>

##  <u>**Introduction**</u> 👋🏻

This is a simple k8s-api project. It is built on **Golang** and utilises the **client-go** library to interact with Kubernetes Cluster.
It is a plug and play solution and can be used to create a k8s-api server. Aim is to build a simple k8s-api server that can be used to interact with any Kubernetes cluster.

 In my industrial experience I have realised that Kuberenetes is a very powerful tool but, only used by a handful of people. Problem is not grasping the concept of the cluster. My last Product Manager was much versed in AWS than I am. 
The reason I feel that Kuberenetes is still only used by a couple of Developers in every firm because there is no easier way to interact with it.
This project will provide a bunch of API endpoints to perform various funstions on the cluster. For now, I will have the Postman collections and API docs to achieve it, plan is build a Dashboard on the API later.

<hr>

## <u>**Getting started**</u> ▶️

To start using the project you need a Kubernetes Cluster and should have access to apply changes to the cluster.
The project can be run *inside* the cluster as well as from *outside* the cluster. We will first discuss how to run the project from *inside* the cluster.
### **Steps to follow** 🪜
- ## **Service Account**:

  We need to make a custom service account to be able to interact with the cluster. We will use this service account in our Pod on which we will run the API.</br>
 Command to make a service account: </br>
    ```
    kubectl apply -f - <<EOF
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: <Your-Custom-Name>
    EOF
    ```
  Or you can apply the YAML file using:
  ```
  kubectl apply -f https://raw.githubusercontent.com/kitarp29/k8s-api/master/yamls/sa.yaml
  ```
   > We can also use the *default* service account. But it is not recommended.
  
  > Verify the service account by running the following command: kubectl get serviceaccounts/Your-Custom-Name -o yaml

  <br>
- ## **Cluster Role**:
  
  We need to make a custom cluster role to be able to interact with the cluster. We will use this cluster role to bind to our **Service Account**. The role should have permission to all the resources in order for the project to run smoothly.
  
  I would advice not to make any Role and use the *cluster-admin* role directly. Still if you want to create a custom role, you can do so. [Here](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
  </br>

- ## **Cluster Role Binding**:

  We will bind the **Service Account** to the **Cluster Role** we just created. To do so the commands needed are:

  ```

  kubectl apply -f - <<EOF
  apiVersion: rbac.authorization.k8s.io/v1
  kind: ClusterRoleBinding
  metadata:
    name: <your-custom-name>
  subjects:
  - kind: ServiceAccount
    name: <Name-of-your-serviceaccount> # name of your service account
    namespace: dev # this is the namespace your service account is in
  roleRef: # referring to your ClusterRole
    kind: ClusterRole
    name: cluster-admin # or the custom role you created in the last step
    apiGroup: rbac.authorization.k8s.io
    EOF
  ```

  Or you can apply the YAML file using:

  ```
  kubectl apply -f https://raw.githubusercontent.com/kitarp29/k8s-api/master/yamls/crb.yaml
  ```

- ## **Deploying the Pod**

  This is it! Now that we have the service account and the cluster role binding, we can deploy the pod. We can use this commad to deploy the pod:

  ```
  kubectl apply -f - <<EOF 
  apiVersion: v1
  kind: Pod
  metadata:
    name: <your-custom-name>
  spec:
    serviceAccount: <Name-of-your-serviceaccount> # name of your service account
    containers:
    - name: <your-custom-name>
      image: kitarp29/k8s-api:2.0
      ports:
      - containerPort: 8000
  EOF
  ```
  or you can apply the YAML file using:

  ```
  kubectl apply -f https://raw.githubusercontent.com/kitarp29/k8s-api/master/yamls/pod.yaml
  ```
- ## **Port-forward**:
  
  Now that we have the pod deployed, we can use the *port-forward* command to access the pod. Expose a port of the pod to the local machine.
  To do so, we can use the following command:

  ```
  kubectl port-forward <your-pod-name> 8000:8080
  ```
> API will be running at <a href="localhost:8000/"> localohost:8000/ </a> now.
  ## **The project is up and running🔥!**
  <hr>

##  <u>**Project Features**</u> 🤯
  -  Get Details about any resource in the cluster.
  -  Create new resources in the cluster.
  -  Delete resources in the cluster.
  -  Run CLI Commands using the API.
  -  Manage HELM Charts.
  -  You can add, install, delete and update HELM charts.
  -  Get live events from the cluster.
  -  It is a REST API to interact with the cluster.
  -  It has health check endpoint as well.
  -  More coming soon... 🚧

<hr>

##  <u>**API Docs**</u> 📖

  There are multiple endpoints in the API. You can find all the endpoints in the [API Docs](https://github.com/kitarp29/k8s-api/master/API_DOCS.md)

  Moreover you can find the **Postman Collections** [Here]() 📮

  <hr>

  ## <u>**Contributors Guide**</u> 🥰
  
  Thanks for considering to contribute to the project. If you have any questions, please feel free to contact me at [Twiiter](https://twitter.com/kitarp29).
  The Contributors Guide is available [Here](https://github.com/kitarp29/k8s-api/master/CONTRIBUTING.md) 📖

  <hr>

  ## <u>**License**</u> 🍻

  This project is licensed under the **MIT license**. Feel free to use it and if you want to contribute, please feel free to fork the project and make a pull request. Thanks!

  <hr>

  ## <u>**FAQ**</u> 🤔

  - **Is this a Unique Product?**
  
      No, this is not a unique product. There are similar impelementaion made by other people.
  
  - **Purpose of this project?**

     It's a pet project to learn *Kubernetes* and *Golang*. I wanted to build this to get a better understanding of these two technologies. I also explored *Docker*.



### Thanks for the interest in my API :)
<hr>