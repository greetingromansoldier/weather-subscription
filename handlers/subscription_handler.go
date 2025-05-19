package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"weather-subscription/models"
	"weather-subscription/storage"
)

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var subRequest models.Subscription
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		log.Printf("Received Content-Type: %s. Assuming JSON for now.", contentType)
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&subRequest); err != nil {
		log.Printf("Error decoding subscription request: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if subRequest.Email == "" || subRequest.City == "" || subRequest.Frequency == "" {
		http.Error(w, "Invalid input: email, city, and frequency are required", http.StatusBadRequest)
		return
	}
	if subRequest.Frequency != models.Hourly && subRequest.Frequency != models.Daily {
		http.Error(w, "Invalid input: frequency must be 'hourly' or 'daily'", http.StatusBadRequest)
		return
	}
	newSubscription := models.Subscription{
		Email:     subRequest.Email,
		City:      subRequest.City,
		Frequency: subRequest.Frequency,
		Confirmed: false,
	}

	confirmationToken := fmt.Sprintf("%s-%d", strings.ReplaceAll(subRequest.Email, "@", "-at-"), time.Now().UnixNano())

	err := storage.StorePendingSubscription(newSubscription, confirmationToken)
	if err != nil {
		if strings.Contains(err.Error(), "email already subscribed and confirmed") {
			log.Printf("Attempt to subscribe already confirmed email: %s", newSubscription.Email)
			http.Error(w, "Email already subscribed and confirmed", http.StatusConflict) // 409
		} else {
			log.Printf("Error storing pending subscription for %s: %v", newSubscription.Email, err)
			http.Error(w, "Failed to process subscription", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Pending subscription stored for %s with token %s.", newSubscription.Email, confirmationToken)

	confirmationLink := fmt.Sprintf("http://localhost:%s/confirm/%s", "8080", confirmationToken)
	log.Printf("SIMULATING SENDING EMAIL to %s: Please confirm your subscription by visiting %s", newSubscription.Email, confirmationLink)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseMessage := map[string]string{"message": "Subscription successful. Confirmation email sent."}
	json.NewEncoder(w).Encode(responseMessage)
}

func ConfirmSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/confirm/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		http.Error(w, "Invalid token: token is missing in path", http.StatusBadRequest)
		return
	}
	token := pathParts[0]
	log.Printf("Received confirmation request for token: %s", token)

	pendingSub, err := storage.FindPendingSubscriptionByToken(token)
	if err != nil || pendingSub == nil {
		log.Printf("Error finding pending subscription by token %s: %v", token, err)
		http.Error(w, "Token not found, invalid, or subscription already confirmed", http.StatusNotFound) // 404
		return
	}

	activeSub, err := storage.FindActiveSubscriptionByEmail(pendingSub.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Database error checking active subscription for %s: %v", pendingSub.Email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if activeSub != nil && activeSub.Confirmed {
		log.Printf("Subscription for %s already confirmed.", pendingSub.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseMessage := map[string]string{"message": "Subscription is already confirmed"}
		json.NewEncoder(w).Encode(responseMessage)
		return
	}

	unsubscribeToken, err := storage.ConfirmSubscriptionByEmailAndToken(pendingSub.Email, token)
	if err != nil {
		log.Printf("Error confirming subscription for email %s with token %s: %v", pendingSub.Email, token, err)
		http.Error(w, "Failed to confirm subscription. Token may be invalid or already used.", http.StatusBadRequest)
		return
	}

	log.Printf("Subscription confirmed for email: %s. Unsubscribe token: %s", pendingSub.Email, unsubscribeToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseMessage := map[string]string{"message": "Subscription confirmed successfully"}
	json.NewEncoder(w).Encode(responseMessage)
}

func UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/unsubscribe/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		http.Error(w, "Invalid unsubscribe token: token is missing in path", http.StatusBadRequest)
		return
	}
	unsubscribeToken := pathParts[0]
	log.Printf("Received unsubscribe request for token: %s", unsubscribeToken)

	err := storage.DeleteSubscriptionByUnsubscribeToken(unsubscribeToken)
	if err != nil {
		log.Printf("Error unsubscribing with token %s: %v", unsubscribeToken, err)
		http.Error(w, "Failed to unsubscribe. Token may be invalid or subscription not found.", http.StatusNotFound)
		return
	}

	log.Printf("Successfully unsubscribed using token: %s", unsubscribeToken)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseMessage := map[string]string{"message": "Unsubscribed successfully"}
	json.NewEncoder(w).Encode(responseMessage)
}
