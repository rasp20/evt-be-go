package main

import (
	"evt-be-go/configs"
	"evt-be-go/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	//run database
	configs.ConnectDB()

	//routes
	routes.EventRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}

// Reference: https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-echo-version-2gdg
