# **API Documentation**

This document is a collection of API endpoints for this repository. It will explain the endpoints and how to use them. The endpoints are grouped into sections based on the path. The method, parmaters and response are described in below.

### **Postman collection** : [Here]()

<hr>

## Kuberenetes Management Routes:

- **Home**
    ```
    Method: GET
    Endpoint: /
    Parametes: None
    Response:
        - httpStatusOk: 200
        - message: Yes! I am alive!
        - type: string
    ```
- **Pods**
    ```
    Method: GET
    Endpoint: /pods
    Parametes:
        - namespace: <namespace>
        - containerDetails: <true/false>
    Response:
        - httpStatusOk: 200
        - message: List of pods
        - type: array
    ```
- **Namespace**
  ```
    Method: GET
    Endpoint: /namespace
    Parametes: None
    Response:
        - httpStatusOk: 200
        - message: List of namespaces
        - type: array
    ```
- **Deployments**
    ```
    Method: GET
    Endpoint: /deployments
    Parametes:
        - namespace: <namespace>
    Response:
        - httpStatusOk: 200
        - message: List of deployments
        - type: array
    ```
- **ConfigMaps**
    ```
    Method: GET
    Endpoint: /configmaps
    Parametes:
        - namespace: <namespace>
    Response:
        - httpStatusOk: 200
        - message: List of configmaps
        - type: array
    ```
- **Services**
    ```
    Method: GET
    Endpoint: /services
    Parametes:
        - namespace: <namespace>
    Response:
        - httpStatusOk: 200
        - message: List of services
        - type: array
    ```
- **Events**
    ```
    Method: GET
    Endpoint: /events
    Parametes:
        - namespace: <namespace>
    Response:
        - httpStatusOk: 200
        - message: List of events
        - type: array
    ```
- **Secrets**
    ```
    Method: GET
    Endpoint: /secrets
    Parametes:
        - namespace: <namespace>
    Response:
        - httpStatusOk: 200
        - message: List of secrets
        - type: array
    ```
- **ReplicationControllers**
    ```
    Method: GET
    Endpoint: /replicationcontrollers
    Parametes:
        - namespace: <namespace>
    Response:
        - httpStatusOk: 200
        - message: List of replicationcontrollers
        - type: array
    ```
- **DaemonSets**
    ```
    Method: GET
    Endpoint: /daemonsets
    Parametes:
        - namespace: <namespace>
    Response:
        - httpStatusOk: 200
        - message: List of daemonsets
        - type: array
    ```
- **Pod Logs**
    ```
    Method: GET
    Endpoint: /podlogs
    Parametes:
        - namespace: <namespace>
        - pod: <pod>
    Response:
        - httpStatusOk: 200
        - message: Pod logs
        - type: string
    ```
- **Create Namespace**
    ```
    Method: POST
    Endpoint: /createNamespace
    Parametes:
        - namespace: <namespace>
    Response:
        - httpStatusOk: 200
        - message: Namespace created
        - type: string
    ```
- **Delete Namespace**
    ```
    Method: DELETE
    Endpoint: /deleteNamespace
    Parametes:
        - namespace: <namespace>
    Response:
        - httpStatusOk: 200
        - message: Namespace deleted
        - type: string
    ```
- **Delete Deployment**
    ```
    Method: DELETE
    Endpoint: /deleteDeployment
    Parametes:
        - namespace: <namespace>
        - deployment: <deployment>
    Response:
        - httpStatusOk: 200
        - message: Deployment deleted
        - type: string
    ```
- **Delete Service**
    ```
    Method: DELETE
    Endpoint: /deleteService
    Parametes:
        - namespace: <namespace>
        - service: <service>
    Response:
        - httpStatusOk: 200
        - message: Service deleted
        - type: string
    ```
- **Delete ConfigMap**
    ```
    Method: DELETE
    Endpoint: /deleteConfigMap
    Parametes:
        - namespace: <namespace>
        - configmap: <configmap>
    Response:
        - httpStatusOk: 200
        - message: ConfigMap deleted
        - type: string
    ```
- **Delete Secret**
    ```
    Method: DELETE
    Endpoint: /deleteSecret
    Parametes:
        - namespace: <namespace>
        - secret: <secret>
    Response:
        - httpStatusOk: 200
        - message: Secret deleted
        - type: string
    ```
- **Delete ReplicationController**
    ```
    Method: DELETE
    Endpoint: /deleteReplicationController
    Parametes:
        - namespace: <namespace>
        - replicationcontroller: <replicationcontroller>
    Response:
        - httpStatusOk: 200
        - message: ReplicationController deleted
        - type: string
    ```
- **Delete DaemonSet**
    ```
    Method: DELETE
    Endpoint: /deleteDaemonSet
    Parametes:
        - namespace: <namespace>
        - daemonset: <daemonset>
    Response:
        - httpStatusOk: 200
        - message: DaemonSet deleted
        - type: string
    ```
- **Delete Pod**
    ```
    Method: DELETE
    Endpoint: /deletePod
    Parametes:
        - namespace: <namespace>
        - pod: <pod>
    Response:
        - httpStatusOk: 200
        - message: Pod deleted
        - type: string
    ```
- **Delete Event**
    ```
    Method: DELETE
    Endpoint: /deleteEvent
    Parametes:
        - namespace: <namespace>
        - event: <event>
    Response:
        - httpStatusOk: 200
        - message: Event deleted
        - type: string
    ```

<hr>

## Apply YAML/JSON Files

- **Apply**

    > The Container will no the YAML file preloaded in the container. You can download them using wget or curl.

        ```
        Method: POST
        Endpoint: /apply
        Parametes:
            - filepath: <filepath>
        Response:
            - httpStatusOk: 200
            - message: YAML/JSON file applied
            - type: string
        ```
<hr>

## Help Routes

- **Helm Repo Add**
    ```
    Method: POST
    Endpoint: /helmRepoAdd
    Parametes:
        - url: <url>
        - repoName: <repo>
    Response:
        - httpStatusOk: 200
        - message: Repo added
        - type: string
    ```
- **Helm Repo Udpate**
    ```
    Method: GET
    Endpoint: /helmRepoUpdate
    Parameters: None
    Response:
        - httpStatusOk: 200
        - message: Repo updated
        - type: string
    ```
- **Helm Install**
    ```
    Method: POST
    Endpoint: /helmInstall
    Parametes:
        - chartName: <chart>
        - chartVersion: <version>
        - namespace: <namespace>
        - values: <values>
    Response:
        - httpStatusOk: 200
        - message: Helm installed
        - type: string
    ```
- **Helm Delete**
    ```
    Method: DELETE
    Endpoint: /helmDelete
    Parametes:
        - name : <chart-name>
        - namespace: <namespace>
    Response:
        - httpStatusOk: 200
        - message: Helm deleted
        - type: string
    ```
<hr>

🚧 **More Routes under Construction**👷

Thanks for your pateince! 🥰