package main

import (
	"fmt"
	"log"
	"net/http"
	"weather-subscription/handlers"
)

func main() {
	fmt.Println("Hello world!")

	http.HandleFunc("/weather", handlers.WeatherHandler)

	port := "8080"
	log.Printf("Server is started. The port is %s...", port)
	log.Printf("Available endpoints: ")
	log.Printf(" GET http://localhost:%s/weather?city=CITY_NAME", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("error occured when launching server: ", err)
	}
}
