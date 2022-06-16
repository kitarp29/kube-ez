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

	// Run Server
	e.Logger.Fatal(e.Start(":8000"))
}
