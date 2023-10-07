package handlers

import (
	"car_dealership/internal/manufacturer"
	"encoding/json"
	"net/http"
)

func GetManufacturers(w http.ResponseWriter, r *http.Request) {
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	m, err := manufacturer.Get()
	if err != nil {
		http.Error(w, "Unable to get manufaturers", http.StatusInternalServerError)
		return
	}

	// Encode the manufacturers slice as JSON and send it as the response
	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateManufacturer(w http.ResponseWriter, r *http.Request) {
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON request body into a Car struct.
	var mData manufacturer.Manufacturer
	err := json.NewDecoder(r.Body).Decode(&mData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = manufacturer.Validate(mData.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	newManufacturer := manufacturer.New(mData.Name)
	if newManufacturer == nil {
		http.Error(w, "Failed to create the manufacturer", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newManufacturer)
}
