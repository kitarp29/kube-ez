# **K8s-API**
[![GitHub contributors](https://img.shields.io/github/contributors/kitarp29/k8s-api)](https://github.com/GDSC-DSI/Schedura/graphs/contributors) [![GitHub issues](https://img.shields.io/github/issues/HAC-2020/Aimers)](https://github.com/kitarp29/k8s-api/issues/) 
![GitHub stars](https://img.shields.io/github/stars/kitarp29/k8s-api) ![GitHub releases](https://img.shields.io/github/release/kitarp29/k8s-api)![GitHub license](https://img.shields.io/github/license/kitarp29/k8s-api)

This is a simple k8s-api project. It is built on **Golang** and utilises the **client-go** library to interact with Kubernetes API. It is also built on **Docker**.
<img src="golang.jpeg"></img>

```
kubectl apply -f - <<EOF 
apiVersion: v1
kind: Pod
metadata:
  name: k8s-api
spec:
  serviceAccount: k8s-api
  containers:
  - name: k8s-api
    image: kitarp29/k8s-api:2.0
    ports:
    - containerPort: 8000
EOF
```

