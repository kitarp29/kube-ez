package main

import (
	api "k8-api/api"
	apply "k8-api/apply"
	"k8-api/cmd"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	api.Main()
	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	// Root route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Yes! I am alive!\n")
	})

	//route to get the k8s cluster info
	// e.POST("/kubeconfig", func(c echo.Context) error {
	// 	//kubeconfig := c.FormValue("kubeconfig")
	// 	//fmt.Println(api.Values())
	// 	return c.String(http.StatusOK, "Kubeconfig is paased\n")
	// })

	//route to get the Pods info in a namespace
	e.GET("/pods", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		containerDetails := c.QueryParam("containerDetails") == "True" || c.QueryParam("containerDetails") == "true"
		return c.String(http.StatusOK, api.Pods(namespace, containerDetails))
	})

	e.GET("/namespace", func(c echo.Context) error {

		return c.String(http.StatusOK, api.NameSpace())
	})

	e.GET("/deployments", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		return c.String(http.StatusOK, api.Deployments(namespace))
	})

	e.GET("/configmaps", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		return c.String(http.StatusOK, api.Configmaps(namespace))
	})

	e.GET("/services", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		return c.String(http.StatusOK, api.Services(namespace))
	})

	e.GET("/events", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		return c.String(http.StatusOK, api.Events(namespace))
	})

	e.GET("/secrets", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		return c.String(http.StatusOK, api.Secrets(namespace))
	})

	e.GET("/replicationController", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		return c.String(http.StatusOK, api.ReplicationController(namespace))
	})

	e.GET("/daemonset", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		return c.String(http.StatusOK, api.DaemonSet(namespace))
	})

	e.GET("/podLogs", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		pod := c.QueryParam("pod")
		return c.String(http.StatusOK, api.PodLogs(namespace, pod))
	})

	e.POST("/command", func(c echo.Context) error {
		commands := c.QueryParam("command")
		return c.String(http.StatusOK, cmd.Command(commands))
	})

	e.POST("/createNamespace", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		return c.String(http.StatusOK, api.CreateNamespace(namespace))
	})

	e.POST("/applyFile", func(c echo.Context) error {
		filepath := c.FormValue("filepath")
		return c.String(http.StatusOK, apply.Main(filepath))
	})

	// e.POST("/applyOnlineFile", func(c echo.Context) error {
	// 	filepath := c.FormValue("filepath")
	// 	return c.String(http.StatusOK, api.OnlineDyanmicClient(filepath))
	// })
	e.DELETE("/deleteNamespace", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		return c.String(http.StatusOK, api.DeleteNamespace(namespace))
	})

	e.DELETE("/deleteDeployment", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		deployment := c.FormValue("deployment")
		return c.String(http.StatusOK, api.DeleteDeployment(namespace, deployment))
	})

	e.DELETE("/deleteService", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		service := c.FormValue("service")
		return c.String(http.StatusOK, api.DeleteService(namespace, service))
	})

	e.DELETE("/deleteConfigMap", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		configMap := c.FormValue("configMap")
		return c.String(http.StatusOK, api.DeleteConfigMap(namespace, configMap))
	})

	e.DELETE("/deleteSecret", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		secret := c.FormValue("secret")
		return c.String(http.StatusOK, api.DeleteSecret(namespace, secret))
	})

	e.DELETE("/deleteReplicationController", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		replicationController := c.FormValue("replicationController")
		return c.String(http.StatusOK, api.DeleteReplicationController(namespace, replicationController))
	})

	e.DELETE("/deleteDaemonSet", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		daemonSet := c.FormValue("daemonSet")
		return c.String(http.StatusOK, api.DeleteDaemonSet(namespace, daemonSet))
	})

	e.DELETE("/deletePod", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		pod := c.FormValue("pod")
		return c.String(http.StatusOK, api.DeletePod(namespace, pod))
	})

	e.DELETE("/deleteEvent", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		event := c.FormValue("event")
		return c.String(http.StatusOK, api.DeleteEvent(namespace, event))
	})
	// Run Server
	e.Logger.Fatal(e.Start(":8000"))
}
