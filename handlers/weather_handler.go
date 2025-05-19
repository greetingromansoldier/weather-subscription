package handlers

import (
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	// "weather-subscription/models"
	"weather-subscription/services"
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")

	if city == "" {
		http.Error(w, "Please, choose 'city' parameter", http.StatusBadRequest)
		return
	}

	weatherData, err := services.GetWeatherDataFromAPI(city)

	if err != nil {
		log.Printf("Error getting weather data from API for city %s: %v", city, err)
		http.Error(w, "Failed when retrieving weather data", http.StatusInternalServerError)
		return
	}

	// if we here, weatherData isnt nil and err is nil
	// in other words, we are alive
	jsonData, marshalErr := json.Marshal(weatherData)
	if marshalErr != nil {
		http.Error(w, "error preparing response data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(jsonData)
	if writeErr != nil {
		log.Printf("Error writing response to client: %v", writeErr)
	}

}
