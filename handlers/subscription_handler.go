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
var pendingConfirmations = make(map[string]models.Subscription)

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

	pendingConfirmations[token] = newSubscription

	// response for client
	log.Printf("Subscription pending for %s with token %s. Details: %+v", newSubscription.Email, token, newSubscription)

	confirmationLink := fmt.Sprintf("http://localhost:%s/confirm/%s", "8080", token)
	log.Printf("SIMULATING SENDING EMAIL to %s: Please confirm your subscription by visiting %s", newSubscription.Email, confirmationLink)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseMessage := map[string]string{"message": "Subscription successful. Confirmation email sent."}
	json.NewEncoder(w).Encode(responseMessage)

	log.Printf("Current pending confirmations: %+v", pendingConfirmations)
	log.Printf("Current active subscriptions: %+v", activeSubscriptions)
}

func ConfirmSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// get token from path

	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/confirm/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		http.Error(w, "Invalid token: token is missing in path", http.StatusBadRequest)
		return
	}

	token := pathParts[0]
	log.Printf("Received confirmation request for token: %s", token)

	subscriptionToConfirm, found := pendingConfirmations[token]
	if !found {
		log.Printf("Token not found or already used: %s", token)
		http.Error(w, "Token not found or invalid", http.StatusNotFound)
		return
	}

	if sub, isActive := activeSubscriptions[subscriptionToConfirm.Email]; isActive && sub.Confirmed {
		log.Printf("Subscription for %s already confirmed.", subscriptionToConfirm.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseMessage := map[string]string{"message": "Subscription is ralready confirmed"}
		json.NewEncoder(w).Encode(responseMessage)
		delete(pendingConfirmations, token)
		return
	}

	subscriptionToConfirm.Confirmed = true

	activeSubscriptions[subscriptionToConfirm.Email] = subscriptionToConfirm

	delete(pendingConfirmations, token)

	log.Printf("Subscription confirmed for email: %s. Details: %+v", subscriptionToConfirm.Email, subscriptionToConfirm)
	log.Printf("Current pending confirmations: %+v", pendingConfirmations)
	log.Printf("Current active subscriptions: %+v", activeSubscriptions)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseMessage := map[string]string{"message": "Subscription confirmed successfully"}
	json.NewEncoder(w).Encode(responseMessage)

}
