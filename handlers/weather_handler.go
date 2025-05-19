package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"weather-subscription/models"
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")

	if city == "" {
		http.Error(w, "Please, choose 'city' parameter", http.StatusBadRequest)
		return
	}

	// temporary code

	var weatherData models.Weather
	var err error

	if city == "London" {
		weatherData = models.Weather{
			Temperature: 15.5,
			Humidity:    12,
			Description: "Today weather is good for a walk",
		}
	} else if city == "Paris" {
		weatherData = models.Weather{
			Temperature: 15.5,
			Humidity:    13,
			Description: "Today weather is good for a walk",
		}
	} else {
		err = fmt.Errorf("weather for city '%s' has not been found", city)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("error when trying get weather data: %v", err), http.StatusInternalServerError)
		return
	}

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
