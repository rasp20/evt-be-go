package routes

import (
	"evt-be-go/controllers"

	"github.com/labstack/echo/v4"
)

func EventRoutes(e *echo.Echo) {
	e.POST("/api/event", controllers.CreateEvent)
	e.GET("/api/event/:eventId", controllers.GetAEvent)
	e.PUT("/api/event/:eventId", controllers.EditAEvent)
	e.DELETE("/api/event/:eventId", controllers.DeleteAEvent)
	e.GET("/api/event", controllers.GetAllEvent)
}
