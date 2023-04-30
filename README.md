# **Kube-ez**
[![GitHub contributors](https://img.shields.io/github/contributors/kitarp29/kube-ez)](https://github.com/kitarp29/kube-ez/graphs/contributors) [![GitHub issues](https://img.shields.io/github/issues/HAC-2020/Aimers)](https://github.com/kitarp29/kube-ez/issues/) 
[![Docker](https://img.shields.io/docker/pulls/kitarp29/k8s-api)](https://hub.docker.com/repository/docker/kitarp29/k8s-api)
[![Docker](https://img.shields.io/docker/stars/kitarp29/k8s-api)](https://hub.docker.com/repository/docker/kitarp29/k8s-api)
![Follow](https://img.shields.io/twitter/follow/kitarp29?label=Follow+Kitarp29&style=social)

<img src="golang.jpeg"></img>

#  <u>**Introduction**</u> 👋🏻

It is built on **Golang** and utilizes the **client-go** library to interact with Kubernetes Cluster.
It is a plug-and-play solution and can be used to create a kube-ez server. In three easy steps, you will have a simple API to interact with your cluster.
The aim is to build a simple API server that can be used to interact with any Kubernetes cluster.

 In my industrial experience, I have realized that **Kubernetes** is a very powerful tool but, only used by a handful of developers in the organization. The problem is not grasping the concept of the cluster. My last Product Manager was much more versed in AWS than I am. 
I feel the reason is that there is no easier way to interact with it.
This project will provide a bunch of API endpoints to perform various functions on the cluster. For now, I will have the Postman collections and API docs to achieve it, plan is to build a Dashboard on the API later.

**Docker Image: [kitarp29/k8s-api](https://hub.docker.com/repository/docker/kitarp29/k8s-api)**
> Use the Docker image with tag 2.0 or above to run the kube-ez server.
<hr>

# <u>**Getting started**</u> ▶️

To start using the project you need a Kubernetes Cluster and should have the right access to apply changes to the cluster.
There are two ways to run this project.
 ## 1. **Outside the Cluster**
 ## 2. **Inside the Cluster**

Steps to run the project are mentioned in the [INSTALL.md](https://github.com/kitarp29/kube-ez/blob/main/INSTALL.md)


  <hr>

#  <u>**Project Features**</u> 🤯
  -  Get details about any resource in the cluster.
  -  It *detects* if you are trying to run the project inside or outside of a cluster.
  -  Create new resources in the cluster.
  -  Delete resources in the cluster.
  -  Run CLI commands using the API.
  -  Manage Helm Charts.
  -  You can add, install, delete and update HELM charts.
  -  Get live events from the cluster.
  -  It is a REST API to interact with the cluster.
  -  It has a health check endpoint as well.
  -  More coming soon... 🚧

<hr>

#  <u>**API Docs**</u> 📖

  There are multiple endpoints in the API. You can find all the endpoints in the [API Docs](https://github.com/kitarp29/kube-ez/blob/main/API_DOCS.md)

  Moreover you can find the **Postman Collections** [Here](https://www.getpostman.com/collections/b14cdaad336ab81340b5) 📮

  <hr>

  # <u>**Contributors Guide**</u> 🥰
  
 Thanks for considering contributing to the project. If you have any questions, please feel free to contact me at [Twitter](https://twitter.com/kitarp29).
  The Contributors Guide is available [Here](https://github.com/kitarp29/kube-ez/blob/main/CONTRIBUTING.md) 📖

  <hr>

  # <u>**License**</u> 🍻

  This project is licensed under the **MIT license**. Feel free to use it and if you want to contribute, please feel free to fork the project and make a pull request. Thanks!

  <hr>

  # <u>**FAQ**</u> 🤔

  - **Is this a Unique Product?**
  
      No, this is not a unique product. There are similar implementations made by other developers.
  
  - **Purpose of this project?**

     It's a pet project to learn *Kubernetes* and *Golang*. I wanted to build this to better understand these two technologies. I also explored *Docker*.



### Thanks for your interest in my API :)
<hr>
