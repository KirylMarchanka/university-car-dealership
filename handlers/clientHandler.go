package handlers

import (
	"car_dealership/internal/client"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetClients(w http.ResponseWriter, r *http.Request) {
	clients := client.GetClients()
	if clients == nil {
		http.Error(w, "Unable to get Clients", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&clients)
}

func GetClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["id"]
	cId, err := strconv.ParseInt(clientID, 10, 64)
	if err != nil {
		http.Error(w, "Incorrect Client id", http.StatusUnprocessableEntity)
	}

	clientData := client.GetClient(cId)
	if clientData == nil {
		http.Error(w, "Unable to get Client", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&clientData)
}

func CreateClient(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a Car struct.
	var clientData client.Client
	err := json.NewDecoder(r.Body).Decode(&clientData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the car data.
	if err := clientData.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Create a new car using the New function.
	newClient := client.New(
		clientData.Name,
		clientData.Phone,
		clientData.Password,
	)

	if newClient == nil || newClient.Id == 0 {
		http.Error(w, "Failed to create the client", http.StatusInternalServerError)
		return
	}

	// Respond with the created car details in JSON format.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Id    int64  `json:"id"`
		Phone string `json:"phone"`
		Name  string `json:"name"`
	}{
		Id:    newClient.Id,
		Phone: newClient.Phone,
		Name:  newClient.Name,
	})
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	c := r.Context().Value("client").(*client.Client)

	if !c.Delete() {
		http.Error(w, "Unable to delete Client", http.StatusInternalServerError)
		return
	}

	// Respond with a success message.
	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content indicates a successful deletion.
}
