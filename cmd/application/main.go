package main

import (
	"log"
	"net/http"
	"os"
	"weather-subscription/handlers"
	"weather-subscription/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // дефолтний порт
	}

	if err := storage.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	if storage.DB != nil {
		defer storage.DB.Close()
	}

	http.HandleFunc("/weather", handlers.WeatherHandler)
	http.HandleFunc("/subscribe", handlers.SubscribeHandler)
	http.HandleFunc("/confirm/", handlers.ConfirmSubscriptionHandler)
	http.HandleFunc("/unsubscribe/", handlers.UnsubscribeHandler)

	log.Printf("Server started on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf(" GET http://localhost:%s/weather?city=CITY_NAME", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error occurred when launching server: %v", err)
	}
}
