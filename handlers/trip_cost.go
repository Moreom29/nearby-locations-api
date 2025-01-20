package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"nearby-locations-api/models"
	"nearby-locations-api/utils"
)

func CalculateTripCost(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Get location ID from URL
	locationIDStr := r.URL.Query().Get("location_id")
	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	// Decode user's current location from request body
	var userLocation struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	err = json.NewDecoder(r.Body).Decode(&userLocation)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch destination location from database
	var destination models.Location
	query := `SELECT id, name, address, latitude, longitude FROM locations WHERE id = $1`
	err = utils.DB.QueryRow(query, locationID).Scan(&destination.ID, &destination.Name, &destination.Address, &destination.Latitude, &destination.Longitude)
	if err != nil {
		http.Error(w, "Destination not found", http.StatusNotFound)
		return
	}

	// Prepare data for TollGuru API
	tollGuruAPIKey := os.Getenv("T6bTrrPnRFM8JD3GD8ngdrTPLRtm3FbM")
	if tollGuruAPIKey == "" {
		http.Error(w, "TollGuru API key not set", http.StatusInternalServerError)
		return
	}

	tollRequest := map[string]interface{}{
		"source":      map[string]float64{"lat": userLocation.Latitude, "lng": userLocation.Longitude},
		"destination": map[string]float64{"lat": destination.Latitude, "lng": destination.Longitude},
	}
	tollRequestJSON, _ := json.Marshal(tollRequest)

	tollURL := "https://apis.tollguru.com/v1/calc/route"
	req, _ := http.NewRequest("POST", tollURL, bytes.NewBuffer(tollRequestJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", tollGuruAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to calculate trip cost", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse TollGuru response
	var tollResponse map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &tollResponse)

	// Extract fuel cost and toll cost
	fuelCost := tollResponse["fuel_cost"].(float64)
	tollCost := tollResponse["toll_cost"].(float64)

	// Respond with total cost
	response := map[string]interface{}{
		"total_cost": fuelCost + tollCost,
		"fuel_cost":  fuelCost,
		"toll_cost":  tollCost,
		"time_ns":    time.Since(start).Nanoseconds(),
	}
	json.NewEncoder(w).Encode(response)
}
