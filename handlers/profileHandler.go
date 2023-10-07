package handlers

import (
	"car_dealership/internal/employee"
	_ "car_dealership/internal/employee"
	"encoding/json"
	"net/http"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	empl := r.Context().Value("employee").(*employee.Employee) // Type assertion to *employee.Employee
	response := map[string]string{"name": empl.Name, "email": empl.Email}
	json.NewEncoder(w).Encode(response)
}
