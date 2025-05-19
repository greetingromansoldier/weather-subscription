package handlers

import (
	"application/models"
	"encoding/json"
	"fmt"
	"net/http"
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
			city:        "London",
			temperature: 15.5,
			description: "Today weather is good for a walk",
		}
	} else if city == "Paris" {
		weatherData = models.Weather{
			city:        "London",
			temperature: 15.5,
			description: "Today weather is good for a walk",
		}
	} else {
		err = fmt.Errorf("weather for city '%s' has not been found", city)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("error when trying get weather data: %v", err), http.StatusInternalServerError)
		return
	}

}
