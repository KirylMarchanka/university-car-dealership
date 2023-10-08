package handlers

import (
	"car_dealership/internal/auth"
	"car_dealership/internal/client"
	"car_dealership/internal/employee"
	"car_dealership/internal/hash"
	"encoding/json"
	"net/http"
)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type clientCredentials struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/Json")

	token := r.Header.Get("Authorization")
	if token != "" {
		http.Error(w, "Already logged in", http.StatusConflict)
		return
	}

	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	empl := employee.Find(creds.Email)
	if empl == nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	if !hash.CheckHash(creds.Password, empl.Password) {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	token, err = auth.GenerateToken(creds.Email)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}

func ClientLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/Json")

	token := r.Header.Get("Authorization")
	if token != "" {
		http.Error(w, "Already logged in", http.StatusConflict)
		return
	}

	var creds clientCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	c := client.GetByPhone(creds.Phone)
	if c == nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	if !hash.CheckHash(creds.Password, c.Password) {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	token, err = auth.GenerateToken(c.Phone)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}
