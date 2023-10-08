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
	authRouter := mux.NewRouter()

	// Apply the middleware to the sub-router
	authRouter.Use(middleware.AuthTokenMiddleware)

	// Define routes within the sub-router that need the middleware
	authRouter.HandleFunc("/cars", handlers.CreateCar).Methods("POST")
	authRouter.HandleFunc("/cars/{id:[0-9]+}", handlers.UpdateCar).Methods("PUT")
	authRouter.HandleFunc("/cars/{id:[0-9]+}", handlers.DeleteCar).Methods("DELETE")

	authRouter.HandleFunc("/manufacturers", handlers.CreateManufacturer).Methods("POST")

	authRouter.HandleFunc("/employees", handlers.GetEmployers).Methods("GET")
	authRouter.HandleFunc("/employees", handlers.CreateEmployee).Methods("POST")
	authRouter.HandleFunc("/employees/{id:[0-9]+}", handlers.GetEmployee).Methods("GET")
	authRouter.HandleFunc("/employees/{id:[0-9]+}", handlers.DeleteEmployee).Methods("DELETE")

	authRouter.HandleFunc("/employees/sale", handlers.CreateNewSale).Methods("POST")

	authRouter.HandleFunc("/clients", handlers.GetClients).Methods("GET")
	authRouter.HandleFunc("/clients/{id:[0-9]+}", handlers.GetClient).Methods("GET")

	authRouter.HandleFunc("/profile", handlers.GetProfile).Methods("GET")

	clientRouter := mux.NewRouter()
	clientRouter.Use(middleware.AuthClientTokenMiddleware)

	r.HandleFunc("/clients", handlers.CreateClient).Methods("POST")
	clientRouter.HandleFunc("/{id:[0-9]+}", handlers.DeleteClient).Methods("DELETE")

	// Define routes outside the sub-router (without middleware)
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/clients/login", handlers.ClientLogin).Methods("POST")
	r.HandleFunc("/cars", handlers.GetCars).Methods("GET")
	r.HandleFunc("/cars/{id:[0-9]+}", handlers.GetCar).Methods("GET")
	r.HandleFunc("/manufacturers", handlers.GetManufacturers).Methods("GET")

	// Merge the sub-router with the main router
	r.PathPrefix("/auth").Handler(http.StripPrefix("/auth", authRouter))
	r.PathPrefix("/clients").Handler(http.StripPrefix("/clients", clientRouter))
}
