package main

import (
	"fmt"
	"log"
	"net/http"
	"weather-subscription/handlers"

	"github.com/joho/godotenv"
)

func main() {
	//env for api key
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file. Make sure it exists in case of running locally")
	}

	fmt.Println("Hello world!")

	//routing
	http.HandleFunc("/weather", handlers.WeatherHandler)
	http.HandleFunc("/subscribe", handlers.SubscribeHandler)
	http.HandleFunc("/confirm/", handlers.ConfirmSubscriptionHandler)

	//server
	port := "8080"
	log.Printf("Server is started. The port is %s...", port)
	log.Printf("Available endpoints: ")
	log.Printf(" GET http://localhost:%s/weather?city=CITY_NAME", port)

	serverErr := http.ListenAndServe(":"+port, nil)
	if serverErr != nil {
		log.Fatal("error occured when launching server: ", serverErr)
	}
}
