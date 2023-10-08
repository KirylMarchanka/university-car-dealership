package handlers

import (
	"car_dealership/internal/car"
	"car_dealership/internal/employee"
	"car_dealership/internal/order"
	"car_dealership/internal/sale"
	"encoding/json"
	"net/http"
	"time"
)

type employeeSale struct {
	CarId   int64 `json:"car_id"`
	OrderId int64 `json:"order_id"`
	Price   int64 `json:"price"`
}

func CreateNewSale(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a Car struct.
	var employeeSale employeeSale
	err := json.NewDecoder(r.Body).Decode(&employeeSale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err = car.Find(employeeSale.CarId); err != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	if !order.Exists(employeeSale.OrderId) {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	empl := r.Context().Value("employee").(*employee.Employee) // Type assertion to *employee.Employee
	if sale.Exists(empl.Id, employeeSale.OrderId) {
		http.Error(w, "Order already have sale", http.StatusNotFound)
		return
	}
	empl.Sale(employeeSale.CarId, employeeSale.OrderId, employeeSale.Price, time.Now().Format("2006-01-02"))

	// Respond with a success message.
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Ok")
}
