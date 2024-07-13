package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/vintedMonitor/types"
)

type DB interface {
	// Define other methods needed
	Get_all_users() (*sql.Rows, error)
	Get_Subscription_of_user(username string) ([]types.Subscription, error)
}
type MyDB struct {
	DB *sql.DB
}

func Connect_to_database() (*MyDB, error) {
	db, err := sql.Open("sqlite3", "vinted.db")
	if err != nil {
		log.Println("Error opening database:", err)
		return nil, err
	}

	// Ping database to ensure connection is established
	err = db.Ping()
	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil, err
	}

	log.Println("Connected to SQLite database")

	return &MyDB{DB: db}, nil
}
func (db *MyDB) Get_all_users() (*[]string, error) {
	rows, err := db.DB.Query("SELECT DISTINCT username FROM options")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		usernames = append(usernames, username)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return &usernames, nil
}

func (db *MyDB) Get_Subscription_of_user(username string) ([]types.Subscription, error) {
	query := "SELECT id, url, webhook FROM options WHERE username = ?"
	rows, err := db.DB.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// List to store subscriptions
	var subscriptions []types.Subscription
	for rows.Next() {
		var subscription types.Subscription

		// Scan row data into struct fields
		err := rows.Scan(&subscription.ID, &subscription.Url, &subscription.Webhook)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Append subscription to the list
		subscriptions = append(subscriptions, subscription)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return subscriptions, nil
}
