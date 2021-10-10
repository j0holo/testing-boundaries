package main

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=user password=password dbname=my_database sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = createSubscribersTable(db)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	app := &application{
		db,
		rdb,
	}

	go func() {
		ticker := time.NewTimer(time.Minute)
		for range ticker.C {
			// The os.Stdout can be replaced with a S3 blob storage API for example
			readEmails(rdb, "subscribers", os.Stdout)
		}
	}()

	router := newRouter(app)

	log.Fatal(http.ListenAndServe(":8080", router))
}
