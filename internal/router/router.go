package router

import (
	"car_dealership/handlers"
	"github.com/gorilla/mux"
)

type errorResponse struct {
	Message string `json:"message"`
}

func Define(r *mux.Router) {
	r.HandleFunc("/cars", handlers.GetCars).Methods("GET")
	r.HandleFunc("/cars/{id:[0-9]+}", handlers.GetCar).Methods("GET")
	r.HandleFunc("/cars", handlers.CreateCar).Methods("POST")
	r.HandleFunc("/cars/{id:[0-9]+}", handlers.UpdateCar).Methods("PUT")
	r.HandleFunc("/cars/{id:[0-9]+}", handlers.DeleteCar).Methods("DELETE")
}
