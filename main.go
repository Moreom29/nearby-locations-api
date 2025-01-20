package main

import (
	"log"
	"net/http"
	"os"

	"nearby-locations-api/handlers"
	"nearby-locations-api/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	utils.ConnectDB()

	// Initialize router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/locations", handlers.CreateLocation).Methods("POST")
	r.HandleFunc("/locations", handlers.GetLocationsByCategory).Methods("GET")
	r.HandleFunc("/search", handlers.SearchLocations).Methods("POST")
	r.HandleFunc("/trip-cost", handlers.CalculateTripCost).Methods("POST")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
