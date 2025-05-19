package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"weather-subscription/models"
)

// temporary save subs here
// key: email, value: subscription
var activeSubscriptions = make(map[string]models.Subscription)

// temporary save unconfirmed subs with tokens
// key: token, value: email (of user)
var pendingConfirmations = make(map[string]string)

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var subRequest models.Subscription

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "applications/json") {
		log.Printf("Received Content-Type: %s. Assuming JSON for now", contentType)
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&subRequest); err != nil {
		log.Printf("Error decoding subscription request: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// some basic validation

	if subRequest.Email == "" || subRequest.City == "" || subRequest.Frequency == "" {
		http.Error(w, "Invalid input: email, city, and frequency are required", http.StatusBadRequest)
		return
	}
	if subRequest.Frequency != models.Hourly && subRequest.Frequency != models.Daily {
		http.Error(w, "Invalid input: frequency must be 'hourly' or 'daily'", http.StatusBadRequest)
		return
	}
	if _, exists := activeSubscriptions[subRequest.Email]; exists {
		log.Printf("Email %s already subscribed and confirmed", subRequest.Email)
		http.Error(w, "Email already subscribed", http.StatusConflict)
		return
	}

	newSubscription := models.Subscription{
		Email:     subRequest.Email,
		City:      subRequest.City,
		Frequency: subRequest.Frequency,
		Confirmed: false,
	}

	// token generation

	token := fmt.Sprintf("%s-%d", strings.ReplaceAll(subRequest.Email, "@", "-at-"), time.Now().UnixNano())

	pendingConfirmations[token] = newSubscription.Email
	log.Printf("Subscription pending for %s with token %s", newSubscription.Email, token)

	confirmationLink := fmt.Sprintf("http://localhost:%s/confrim/%s", "8080", token)
	log.Printf("==SENDNIG EMAI SIMULATION== to %s: confirm your subscrioption visiting %s", newSubscription.Email, confirmationLink)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseMessage := map[string]string{"message": "Subscription successful. Confirmation email sent"}
	json.NewEncoder(w).Encode(responseMessage)

	log.Printf("Current pending confirmations: %+v", pendingConfirmations)
	log.Printf("Current active subscriptions: %+v", activeSubscriptions)

}
