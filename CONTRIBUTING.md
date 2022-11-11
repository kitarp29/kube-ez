# **Contributing**

Thanks for contributing to the project!🥰
I am glad you are here. Hoping that you liked the project.

##  🔥**Purpose of this file:**
 - I will be stating the rules for contributing to the project. 
 - Most important is that you should be able to run the project locally. 
 - And I will try to make it as easy as possible for you. 
 - Also this is for future *Me*, who will definitely forget how and why I did it, what I did 😅

## 🚨**Rules for contributing**
1. **Always** raise an issue before making the pull request.
2. Try to wait for the issue to get assigned to you.
3. Make **signed** and **small** commits in your Pull Request.
4. **Fork** the Repository and then work.
5. Comment your code as much as possible.
6. Try to write modular code, I also tried to avoid *Spaghetti*🍝 code.

## 💁**How to contribute**
1. **Always** raise an Issue before making the pull request.
2. Make **signed** and **small** commits in your pull request.
3. State the work you did in the issue and in the pull request.
4. **Fork** the Repository and then work.

## 💬**How to run the project locally**
1. Run the following command:
    ```
    $ git clone https://github.com/kitarp29/kube-ez.git
    $ cd kube-ez
    ```
2. Make sure you have **Golang** installed on your system.
   
   ``` 
   go --version 
   ```
    > This project when built used ```go version go1.18.3 linux/amd64```
3. If ```go.mod``` file is not present, then run the following command:
   ```
   go mod init kube-ez
   ```
4. Now, run the following command:
   ```
   go mod tidy
   ```
   > This command will pull all the latest packages from the internet.

5. Now, if you have run it **inside the container**, then follow the steps given in [README.md](https://github.com/kitarp29/kube-ez/blob/main/README.md)
6. If you have to run it **outside the container**:
    - Uncomment the code in api/api.go file from line 97 to 119
    - Comment out the code in api/api..go from 123 to 131
    - The comments are self-explaining.
    - Export a variable called ```KUBECONFIG``` and set it to the path of the kube-ez kubeconfig file. (Usually ~/.kube/config)
    - Run the following command:
        ```
        export KUBECONFIG=~/.kube/config
        ```
    - Now, run the following command:
        ```
        go run server.go
        ```
> Reason for these two ways of running the project:
> The Client-go library needs a file that tells it where to find the kubeconfig file. If the code is inside the container, then the file is accessible. Outside the cluster, we need to explicitly export the KUBECONFIG variable.
   

## 🐋 **Docker Image** [Link](https://hub.docker.com/repository/docker/kitarp29/k8s-api)
It's a basic container based on the latest release of **Golang**. The tag *2.0* works well.

## 📂 **File Structure**
1. **api**:
    - **api.go**:
         This file contains the main logic of the project. It has all the functions that interact with the *client-go* library. It also has the ```main()``` function that starts with the server. It will help us run the project even outside the cluster.
2. **install**:
    - **install.go**:
        This file contains the logic of the **install** command. It will apply the changes to the cluster. It handles all the requests related to **Helm** charts. It helps us add/upgrade/delete the charts.
3. **apply**:
    - **apply.go**: 
        This file contains the logic of the **apply** command. It will apply the changes to the cluster. It helps apply any YAML /JSON File to our cluster.
4. **yamls**:
    - **sa.yaml**: YAML to apply desired ServiceAccount for the project.
    - **crb.yaml**: YAML to apply desired CustomResourceDefinition for the project.
    - **pod.yaml**: YAML to apply the desired Pod for the project.
5. **server.go**
    - This file contains the logic of the **server** command. It will start the server. It will start the server and listen on the port ```8000```. It has all the routes for the project.
6. **Dockerfile**
7. Markdown files
8. License file  
   
   <hr>
