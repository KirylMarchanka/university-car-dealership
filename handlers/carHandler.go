package handlers

import (
	"car_dealership/internal/car"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Parse query parameters with updated naming conventions.
	manufacturerId := r.URL.Query().Get("manufacturer_id")

	var mId int64

	if manufacturerId != "" {
		var err error

		mId, err = strconv.ParseInt(manufacturerId, 10, 64)
		if err != nil {
			http.Error(w, "{\"message\":\"manufacturer_id must be an integer\"}", http.StatusUnprocessableEntity)

			return
		}
	}

	name := r.URL.Query().Get("name")
	fuel := r.URL.Query().Get("fuel")
	ob := r.URL.Query().Get("order_by")
	if ob != "" {
		if ob != "max_speed" && ob != "acceleration" {
			http.Error(w, "{\"message\":\"incorrect order_by\"}", http.StatusUnprocessableEntity)

			return
		}
	}

	od := r.URL.Query().Get("order_direction")

	cars, err := car.SelectCars(mId, name, fuel, ob, od)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON containing the list of cars.
	json.NewEncoder(w).Encode(struct {
		Cars []car.Car `json:"cars"`
	}{
		Cars: cars,
	})
}

func GetCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	c, err := car.Find(id)
	if err != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Car car.Car `json:"car"`
	}{
		Car: c,
	})
}

func CreateCar(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a Car struct.
	var carData car.Car
	err := json.NewDecoder(r.Body).Decode(&carData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the car data.
	if err := car.Validate(
		carData.ManufacturerId,
		carData.Name,
		carData.Fuel,
		carData.FuelCapacity,
		carData.Engine,
		carData.EnginePower,
		carData.EngineCapacity,
		carData.MaxSpeed,
		carData.Acceleration,
	); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Create a new car using the New function.
	newCar := car.New(
		carData.ManufacturerId,
		carData.Name,
		carData.Fuel,
		carData.FuelCapacity,
		carData.Engine,
		carData.EnginePower,
		carData.EngineCapacity,
		carData.MaxSpeed,
		carData.Acceleration,
	)

	if newCar == nil {
		http.Error(w, "Failed to create the car", http.StatusInternalServerError)
		return
	}

	// Respond with the created car details in JSON format.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCar)
}

func UpdateCar(w http.ResponseWriter, r *http.Request) {
	// Extract car ID from the URL path using Gorilla Mux.
	vars := mux.Vars(r)
	carID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	// Parse the JSON request body into a Car struct.
	var updatedCar car.Car
	err = json.NewDecoder(r.Body).Decode(&updatedCar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the updated car data.
	if err := car.Validate(
		updatedCar.ManufacturerId,
		updatedCar.Name,
		updatedCar.Fuel,
		updatedCar.FuelCapacity,
		updatedCar.Engine,
		updatedCar.EnginePower,
		updatedCar.EngineCapacity,
		updatedCar.MaxSpeed,
		updatedCar.Acceleration,
	); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	updatedCar.Id = carID
	if err := car.Update(&updatedCar); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message or updated car details in JSON format.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCar)
}

func DeleteCar(w http.ResponseWriter, r *http.Request) {
	// Extract car ID from the URL path using Gorilla Mux.
	vars := mux.Vars(r)
	carID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	c, err := car.Find(carID)
	if err != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	// Perform the delete operation in your database, e.g., using a Delete function.
	// Replace this with your actual database delete logic.
	if err := car.Delete(&c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message.
	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content indicates a successful deletion.
}
