package main

import (
	"log"
	"net/http"
	"weather-subscription/handlers"
	"weather-subscription/storage"

	"github.com/joho/godotenv"
)

func main() {
	//env for api key
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file.")
	}

	// database init
	if err := storage.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	if storage.DB != nil {
		defer storage.DB.Close()
	}

	//routing
	http.HandleFunc("/weather", handlers.WeatherHandler)
	http.HandleFunc("/subscribe", handlers.SubscribeHandler)
	http.HandleFunc("/confirm/", handlers.ConfirmSubscriptionHandler)
	http.HandleFunc("/unsubscribe/", handlers.UnsubscribeHandler)

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
