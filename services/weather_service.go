package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"weather-subscription/models"
)

// weather api key implemented in .env file
const weatherAPIBaseURL = "http://api.weatherapi.com/v1/current.json"

type WeatherAPIResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`

	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  int     `json:"humidity"`
		Condition struct {
			Text string `json:"condition:text"`
		} `json:"condition"`
	} `json:"current"`
}

func GetWeatherDataFromAPI(city string) (*models.Weather, error) {
	apiKey := os.Getenv("WEATHERAPI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WEATHERAPI_KEY environment variable hasnt been set. It should be in .env file")
	}

	encodedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf("%s?key=%s&q=%s&aqi=no", weatherAPIBaseURL, apiKey, encodedCity)

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error occured when making request to WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("WeatherAPI request failed with statud %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body from WeatherAPI: %w", err)
	}
	var apiResponse WeatherAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling WeatherAPI response: %w. Response body: %s", err, string(bodyBytes))
	}

	weatherData := &models.Weather{
		City:        apiResponse.Location.Name,
		Temperature: apiResponse.Current.TempC,
		Humidity:    float64(apiResponse.Current.Humidity),
		Description: apiResponse.Current.Condition.Text,
	}

	return weatherData, nil

}
