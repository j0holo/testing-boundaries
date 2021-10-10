package main

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
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

	router := newRouter(app)

	log.Fatal(http.ListenAndServe(":8080", router))
}
