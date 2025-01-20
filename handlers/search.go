package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"nearby-locations-api/models"
	"nearby-locations-api/utils"
)

func SearchLocations(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var searchParams struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Category  string  `json:"category"`
		RadiusKM  float64 `json:"radius_km"`
	}

	err := json.NewDecoder(r.Body).Decode(&searchParams)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, address, latitude, longitude, category FROM locations WHERE category = $1`
	rows, err := utils.DB.Query(query, searchParams.Category)
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

		distance := haversine(searchParams.Latitude, searchParams.Longitude, location.Latitude, location.Longitude)
		if distance <= searchParams.RadiusKM {
			locations = append(locations, location)
		}
	}

	response := map[string]interface{}{
		"locations": locations,
		"time_ns":   time.Since(start).Nanoseconds(),
	}
	json.NewEncoder(w).Encode(response)
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in KM
	lat1, lon1, lat2, lon2 = toRadians(lat1), toRadians(lon1), toRadians(lat2), toRadians(lon2)

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func toRadians(degree float64) float64 {
	return degree * math.Pi / 180
}
