package handlers

import (
	"car_dealership/internal/employee"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetEmployers(w http.ResponseWriter, r *http.Request) {
	var err error

	saleMinDate := r.URL.Query().Get("sale_min_date")
	var saleMin time.Time
	if saleMinDate != "" {
		saleMin, err = time.Parse("2006-01-02", saleMinDate)
		if err != nil {
			http.Error(w, "Invalid sale_min_date format", http.StatusUnprocessableEntity)
			return
		}
	}

	saleMaxDate := r.URL.Query().Get("sale_max_date")
	var saleMax time.Time
	if saleMaxDate != "" {
		saleMax, err = time.Parse("2006-01-02", saleMaxDate)
		if err != nil {
			http.Error(w, "Invalid sale_max_date format", http.StatusUnprocessableEntity)
			return
		}
	}

	if (!saleMin.IsZero() && !saleMax.IsZero()) && (saleMin.After(saleMax) || saleMax.Before(saleMin)) {
		http.Error(w, "Invalid date", http.StatusUnprocessableEntity)
	}

	priceMin := r.URL.Query().Get("sale_min_price")
	var priceMinConv int64
	if priceMin != "" {
		priceMinConv, err = strconv.ParseInt(priceMin, 10, 64)
		if err != nil {
			http.Error(w, "Invalid sale_min_price format", http.StatusUnprocessableEntity)
			return
		}
	}

	priceMax := r.URL.Query().Get("sale_max_price")
	var priceMaxConv int64
	if priceMax != "" {
		priceMaxConv, err = strconv.ParseInt(priceMax, 10, 64)
		if err != nil {
			http.Error(w, "Invalid sale_max_price format", http.StatusUnprocessableEntity)
			return
		}
	}

	if (priceMaxConv != 0 && priceMinConv != 0) && ((priceMinConv > priceMaxConv) || (priceMaxConv < priceMinConv)) {
		http.Error(w, "Incorrect prices", http.StatusUnprocessableEntity)
		return
	}

	employees, err := employee.GetEmployers(saleMinDate, saleMaxDate, priceMinConv, priceMaxConv)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Failed to get employees", http.StatusInternalServerError)
		return
	}

	// Create a JSON response
	response := map[string][]map[string]interface{}{"employees": employees}

	// Set the Content-Type header and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	empl := employee.GetEmployee(id)
	if empl == nil || empl.Email == "" {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	// Set the Content-Type header and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(empl)
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a Car struct.
	var employeeData employee.Employee
	err := json.NewDecoder(r.Body).Decode(&employeeData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the car data.
	if err := employee.Validate(
		employeeData.Name,
		employeeData.Email,
		employeeData.Password,
	); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Create a new car using the New function.
	newEmployee := employee.New(
		employeeData.Name,
		employeeData.Email,
		employeeData.Password,
	)

	if newEmployee == nil {
		http.Error(w, "Failed to create the employee", http.StatusInternalServerError)
		return
	}

	// Respond with the created car details in JSON format.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Id    int64  `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}{
		Id:    newEmployee.Id,
		Email: newEmployee.Email,
		Name:  newEmployee.Name,
	})
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	empl := employee.FindById(id)
	if empl == nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	if !employee.Delete(empl.Id) {
		http.Error(w, "Failed to delete Employee", http.StatusInternalServerError)
		return
	}

	// Respond with a success message.
	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content indicates a successful deletion.
}
