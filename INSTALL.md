# 1. **Outside** the Cluster
This is the easiest way to run this project. Provided you have the **kubeconfig** file in your local.

The code in looks for the file at the path ```"HOME"/.kube/config```

All you have to do after this is:
- Clone the repo
- Open a terminal in the dir
- Run this: 
 ``` 
 go mod tidy && go run k8-api
 ```
-  API will be running at <a href="localhost:8000/"> localhost:8000 </a> now.

  ### **The project is up and running🔥!**

<br>

# 2. **Inside** the Cluster

<font size="4.5"><b>Run this command ✨: </b></font>

``` 
kubectl apply -f https://raw.githubusercontent.com/kitarp29/kube-ez/main/yamls/one-to-rule-them-all.yaml && sleep 5 && kubectl port-forward kube-ez 8000:8000
```
> API will be running at <a href="localhost:8000/"> localhost:8000 </a> now

> P.S.: I am using the *sleep* command to give the pod some time to start. You can change it if you want.

<font size="4.5"><center><b>---OR---</b></center>

If you want to do a custom installation. Follow the steps below:</font>

-  <font size="3"><b>Service Account:</b></font>

  We need to make a custom service account to be able to interact with the cluster. We will use this service account in our pod on which we will run the API.</br>
  Command to make custom a service account: </br>

  ```
    kubectl apply -f - <<EOF
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: <Your-Custom-Name>
    EOF
   ```

  > We can also use the *default* service account. But it is not recommended.
  
-  <font size="3"><b>Cluster Role:</b></font>
  
  We need to make a custom cluster role to be able to interact with the cluster. We will use this cluster role to bind to our **Service Account**. The role should have permission to all the resources in order for the project to run smoothly.
  I would advise you **not** to make any role and use the *cluster-admin* role directly. Still, if you want to create a custom role, you can do so. [Refer Here](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
  </br>

-  <font size="3"><b>Cluster Role Binding:</b></font>

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

-  <font size="3"><b>Deploying the Pod</b></font>

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
-  <font size="3"><b>Port-forward:</b></font>
  
  Now that we have the pod deployed, we can use the *port-forward* command to access the pod. Expose a port of the pod to the local machine.
  To do so, we can use the following command:

  ```
  kubectl port-forward <your-pod-name> 8000:8000
  ```
> API will be running at <a href="localhost:8000/"> localhost:8000 </a> now.
  ### **The project is up and running🔥!**
 <hr>

  Now that the Project is up and running for you. 
  You can learn how to use the API from the [API Docs](https://github.com/kitarp29/kube-ez/blob/main/API_DOCS.md).

  There are multiple endpoints in the API. You can find all the endpoints in the [API Docs](https://github.com/kitarp29/kube-ez/blob/main/API_DOCS.md)

  Moreover you can find the **Postman Collections** [Here](https://www.getpostman.com/collections/b14cdaad336ab81340b5) 📮
 <hr>