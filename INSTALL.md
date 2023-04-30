# 1. **Outside** the Cluster
This is the easiest way to run this project. Provided you have the **kubeconfig** file in your local.

The code in looks for the file at the path ```"HOME"/.kube/config```

All you have to do after this is:
- Clone the repo
- Open a terminal in the dir
- Run these two commands:

- ### ``` go mod tidy```

- ### ``` go run k8-api```


 > API will be running at <a href="localhost:8000/"> localhost:8000 </a> now.
  ### **The project is up and running🔥!**

I usually have a Kind Cluster running locally that I connet to using this project. So I have kept it this way, I am plaaning to add flags in the run command next. The flags will allow user to specify the path of the kubeconfig file. This might be a good task to take up.


<br>

# 2. **Inside** the Cluster

This might ned some extra prerequistes. But it ensures you can connect to any cluster in the world. I actually use this in my CI, helps e test code on a nuetral environment. 

## **Steps to follow** 🪜
- ### **Service Account**:

  We need to make a custom service account to be able to interact with the cluster. We will use this service account in our pod on which we will run the API.</br>
 Command to make a service account: </br>
    ```
    kubectl apply -f - <<EOF
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: <Your-Custom-Name>
    EOF
    ```
  Or you can apply the [YAML file](https://raw.githubusercontent.com/kitarp29/kube-ez/main/yamls/sa.yaml) using:
  ```
  kubectl apply -f https://raw.githubusercontent.com/kitarp29/kube-ez/main/yamls/sa.yaml
  ```
   > We can also use the *default* service account. But it is not recommended.
  
  > Verify the service account by running the following command: kubectl get serviceaccounts/Your-Custom-Name -o yaml
- ### **Cluster Role**:
  
  We need to make a custom cluster role to be able to interact with the cluster. We will use this cluster role to bind to our **Service Account**. The role should have permission to all the resources in order for the project to run smoothly.
  I would advise you not to make any role and use the *cluster-admin* role directly. Still, if you want to create a custom role, you can do so. [Here](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
  </br>

- ### **Cluster Role Binding**:

  We will bind the **Service Account** to the **Cluster Role** we just created. To do so the commands needed are:

  ```
  kubectl apply -f - <<EOF
  apiVersion: rbac.authorization.k8s.io/v1
  kind: ClusterRoleBinding
  metadata:
    name: <your-custom-name>
  subjects:
  - kind: ServiceAccount
    name: <Name-of-your-serviceaccount> 
    namespace: default
  roleRef:
    kind: ClusterRole
    name: cluster-admin
    apiGroup: rbac.authorization.k8s.io
  EOF
  ```

  Or you can apply the [YAML](https://raw.githubusercontent.com/kitarp29/kube-ez/main/yamls/crb.yaml) file using:

  ```
  kubectl apply -f https://raw.githubusercontent.com/kitarp29/kube-ez/main/yamls/crb.yaml
  ```

- ### **Deploying the Pod**

  This is it! Now that we have the service account and the cluster role binding, we can deploy the pod. We can use this command to deploy the pod:

  ```
  kubectl apply -f - <<EOF 
  apiVersion: v1
  kind: Pod
  metadata:
    name: <your-custom-name>
  spec:
    serviceAccount: <Name-of-your-serviceaccount> 
    containers:
    - name: <your-custom-name>
      image: kitarp29/k8s-api:9.0
      ports:
      - containerPort: 8000
  EOF
  ```
  or you can apply the [YAML](https://raw.githubusercontent.com/kitarp29/kube-ez/main/yamls/pod.yaml) file using:

  ```
  kubectl apply -f https://raw.githubusercontent.com/kitarp29/kube-ez/main/yamls/pod.yaml
  ```
- ### **Port-forward**:
  
  Now that we have the pod deployed, we can use the *port-forward* command to access the pod. Expose a port of the pod to the local machine.
  To do so, we can use the following command:

  ```
  kubectl port-forward <your-pod-name> 8000:8000
  ```
> API will be running at <a href="localhost:8000/"> localhost:8000 </a> now.
  ### **The project is up and running🔥!**

<br> <hr>

  Now that the Project is up and running for you. 
  You can learn how to use the API from the [API Docs](https://github.com/kitarp29/kube-ez/blob/main/API_DOCS.md).

  There are multiple endpoints in the API. You can find all the endpoints in the [API Docs](https://github.com/kitarp29/kube-ez/blob/main/API_DOCS.md)

  Moreover you can find the **Postman Collections** [Here](https://www.getpostman.com/collections/b14cdaad336ab81340b5) 📮
