package router

import (
	"evt-be-go/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/event", middleware.GetAllEvent).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/event", middleware.CreateEvent).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/event/{id}", middleware.UpdateEvent).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteEvent/{id}", middleware.DeleteEvent).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/deleteAllEvent", middleware.DeleteAllEvent).Methods("DELETE", "OPTIONS")
	return router
}
