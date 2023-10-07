package handlers

import (
	"car_dealership/internal/car"
	"car_dealership/internal/employee"
	"encoding/json"
	"net/http"
	"time"
)

type employeeSale struct {
	CarId int64 `json:"car_id"`
	Price int64 `json:"price"`
}

func CreateNewSale(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a Car struct.
	var sale employeeSale
	err := json.NewDecoder(r.Body).Decode(&sale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = car.Find(sale.CarId)
	if err != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	empl := r.Context().Value("employee").(*employee.Employee) // Type assertion to *employee.Employee
	empl.Sale(sale.CarId, sale.Price, time.Now().Format("2006-01-02"))

	// Respond with a success message.
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Ok")
}
