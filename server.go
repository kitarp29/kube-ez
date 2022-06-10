package main

import (
	"fmt"
	api "k8-api/api"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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
	e.POST("/kubeconfig", func(c echo.Context) error {
		kubeconfig := c.FormValue("kubeconfig")
		fmt.Println(api.Values(kubeconfig))
		return c.String(http.StatusOK, "Kubeconfig is paased\n")
	})

	//route to get the Pods info in a namespace
	e.GET("/pods", func(c echo.Context) error {
		namespace := c.QueryParam("namepspace")

		return c.String(http.StatusOK, "Namespace!\n"+api.Pods(namespace))
	})

	// Run Server
	e.Logger.Fatal(e.Start(":8000"))
}
