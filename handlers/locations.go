package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"nearby-locations-api/models"
	"nearby-locations-api/utils"
)

func CreateLocation(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var location models.Location
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO locations (name, address, latitude, longitude, category) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = utils.DB.QueryRow(query, location.Name, location.Address, location.Latitude, location.Longitude, location.Category).Scan(&location.ID)
	if err != nil {
		http.Error(w, "Error inserting location", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":      location.ID,
		"time_ns": time.Since(start).Nanoseconds(),
	}
	json.NewEncoder(w).Encode(response)
}

func GetLocationsByCategory(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, address, latitude, longitude, category FROM locations WHERE category = $1`
	rows, err := utils.DB.Query(query, category)
	if err != nil {
		http.Error(w, "Error fetching locations", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var location models.Location
		err := rows.Scan(&location.ID, &location.Name, &location.Address, &location.Latitude, &location.Longitude, &location.Category)
		if err != nil {
			http.Error(w, "Error scanning location", http.StatusInternalServerError)
			return
		}
		locations = append(locations, location)
	}

	response := map[string]interface{}{
		"locations": locations,
		"time_ns":   time.Since(start).Nanoseconds(),
	}
	json.NewEncoder(w).Encode(response)
}
