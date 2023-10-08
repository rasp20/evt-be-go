package routes

import (
	"evt-be-go/controllers"

	"github.com/labstack/echo/v4"
)

func EventRoutes(e *echo.Echo) {
	e.POST("/api/event", controllers.CreateEvent)
	e.GET("/api/event/:eventId", controllers.GetAnEvent)
	e.PUT("/api/event/:eventId", controllers.EditAnEvent)
	e.DELETE("/api/event/:eventId", controllers.DeleteAnEvent)
	e.GET("/api/event", controllers.GetAllEvent)
}
