package router

import (
	"car_dealership/handlers"
	"car_dealership/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func Define(r *mux.Router) {
	// Create a sub-router for routes that need middleware
	authRouter := mux.NewRouter().StrictSlash(false)

	// Apply the middleware to the sub-router
	authRouter.Use(middleware.AuthTokenMiddleware)

	// Define routes within the sub-router that need the middleware
	authRouter.HandleFunc("/cars", handlers.CreateCar).Methods("POST")
	authRouter.HandleFunc("/cars/{id:[0-9]+}", handlers.UpdateCar).Methods("PUT")
	authRouter.HandleFunc("/cars/{id:[0-9]+}", handlers.DeleteCar).Methods("DELETE")
	authRouter.HandleFunc("/profile", handlers.GetProfile).Methods("GET")

	// Define routes outside the sub-router (without middleware)
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/cars", handlers.GetCars).Methods("GET")
	r.HandleFunc("/cars/{id:[0-9]+}", handlers.GetCar).Methods("GET")

	// Merge the sub-router with the main router
	r.PathPrefix("/auth").Handler(http.StripPrefix("/auth", authRouter))
}
