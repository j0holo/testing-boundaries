package main

import (
	"database/sql"
	"time"
)

func createSubscribersTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS subscribers (
    	id serial primary key,
    	name text,
    	email text unique,
    	membership_type text,
    	uid text,
    	created_at timestamp 
	)`)
	return err
}

func saveNewSubscriber(db *sql.DB, name, email, membershipType, UID string, now time.Time) error {
	_, err := db.Exec(`INSERT INTO subscribers (name, email, membership_type, uid, created_at) VALUES ($1, $2, $3, $4, $5)`, name, email, membershipType, UID, now)
	return err
}
