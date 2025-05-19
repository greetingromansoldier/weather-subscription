package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"weather-subscription/models"
)

var DB *sql.DB

const dbFileName = "subscriptions.db"

func InitDB() error {
	var err error
	_ = os.Remove(dbFileName)

	DB, err = sql.Open("sqlite3", "./"+dbFileName)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	log.Println("Successfully connectted to SQLite db")
	return createTables()
}

func createTables() error {
	createSubscriptionsTableSQL := `CREATE TABLE IF NOT EXISTS subscriptions (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"email" TEXT NOT NULL UNIQUE,
		"city" TEXT NOT NULL,
		"frequency" TEXT NOT NULL,
		"confirmed" BOOLEAN NOT NULL DEFAULT 0,
		"confirmation_token" TEXT UNIQUE, 
		"unsubscribe_token" TEXT UNIQUE, 
		"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	statement, err := DB.Prepare(createSubscriptionsTableSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare create subs table statement: %w", err)
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("failed to execute create subs table statement: %w", err)
	}

	log.Println("Table 'subscriptions' created or already exists.")
	return nil
}

func StorePendingSubscription(sub models.Subscription, confirmationToken string) error {
	var confirmed bool
	err := DB.QueryRow("SELECT confirmed FROM subscriptions WHERE email = ?", sub.Email).Scan(&confirmed)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing subscription: %w", err)
	}
	if err == nil && confirmed {
		return fmt.Errorf("email already subscribed and confirmed")
	}

	stmt, err := DB.Prepare(`
		INSERT INTO subscriptions (email, city, frequency, confirmed, confirmation_token, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(email) DO UPDATE SET
			city = excluded.city,
			frequency = excluded.frequency,
			confirmed = excluded.confirmed,
			confirmation_token = excluded.confirmation_token,
			created_at = excluded.created_at;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert/update pending subscription statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(sub.Email, sub.City, sub.Frequency, false, confirmationToken, time.Now())
	if err != nil {
		return fmt.Errorf("failed to execute insert/update pending subscription: %w", err)
	}
	return nil
}

func FindPendingSubscriptionByToken(confirmationToken string) (*models.Subscription, error) {
	var sub models.Subscription
	err := DB.QueryRow(`
		SELECT email, city, frequency, confirmed 
		FROM subscriptions 
		WHERE confirmation_token = ? AND confirmed = 0`,
		confirmationToken,
	).Scan(&sub.Email, &sub.City, &sub.Frequency, &sub.Confirmed)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pending subscription with token not found or already confirmed")
		}
		return nil, fmt.Errorf("error finding pending subscription by token: %w", err)
	}
	return &sub, nil
}

func ConfirmSubscriptionByEmailAndToken(email string, confirmationToken string) (string, error) {
	unsubscribeToken := fmt.Sprintf("unsub-%s-%d", email, time.Now().UnixNano())

	stmt, err := DB.Prepare(`
		UPDATE subscriptions 
		SET confirmed = 1, confirmation_token = NULL, unsubscribe_token = ?
		WHERE email = ? AND confirmation_token = ? AND confirmed = 0 
	`)
	if err != nil {
		return "", fmt.Errorf("failed to prepare confirm subscription statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(unsubscribeToken, email, confirmationToken)
	if err != nil {
		return "", fmt.Errorf("failed to execute confirm subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return "", fmt.Errorf("no pending subscription found for email and token to confirm, or already confirmed")
	}

	return unsubscribeToken, nil
}

func FindActiveSubscriptionByEmail(email string) (*models.Subscription, error) {
	var sub models.Subscription
	err := DB.QueryRow(
		"SELECT email, city, frequency, confirmed FROM subscriptions WHERE email = ? AND confirmed = 1",
		email,
	).Scan(&sub.Email, &sub.City, &sub.Frequency, &sub.Confirmed)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding active subscription by email: %w", err)
	}
	return &sub, nil
}

func DeleteSubscriptionByUnsubscribeToken(unsubscribeToken string) error {
	stmt, err := DB.Prepare("DELETE FROM subscriptions WHERE unsubscribe_token = ? AND confirmed = 1")
	if err != nil {
		return fmt.Errorf("failed to prepare delete subscription statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(unsubscribeToken)
	if err != nil {
		return fmt.Errorf("failed to execute delete subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected on delete: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no active subscription found for the given unsubscribe token")
	}
	return nil
}
