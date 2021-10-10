package main

import "database/sql"

func createSubscribersTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS subscribers (
    	id serial primary key,
    	name text,
    	email text unique,
    	membership_type text,
    	uid text
	)`)
	return err
}

func saveNewSubscriber(db *sql.DB, name, email, membershipType, UID string) error {
	_, err := db.Exec(`INSERT INTO subscribers (name, email, membership_type, uid) VALUES ($1, $2, $3, $4)`, name, email, membershipType, UID)
	return err
}
