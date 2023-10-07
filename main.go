package main

import (
	"car_dealership/internal/router"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	router.Define(r)

	log.Fatal(http.ListenAndServe(":8080", r))
}
