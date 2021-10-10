package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"net/http"
)

func newRouter(a *application) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/api/v1/json", a.jsonHandler)

	return router
}

type application struct {
	db  *sql.DB
	rdb *redis.Client
}

func (a *application) jsonHandler(w http.ResponseWriter, r *http.Request) {
	// Read input from request
	rawInput, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		writeJSONToClient(w, http.StatusBadRequest, Output{})
		return
	}
	defer r.Body.Close()

	input := Input{}
	err = json.Unmarshal(rawInput, &input)
	if err != nil {
		log.Println(err)
		writeJSONToClient(w, http.StatusBadRequest, Output{})
		return
	}

	// Call the actual business logic
	output, err := addNewSubscriber(a.db, a.rdb, input)
	if err != nil {
		if errors.Is(err, errInvalidMembershipType) {
			writeJSONToClient(w, http.StatusBadRequest, output)
			return
		}
		log.Println(err)
		writeJSONToClient(w, http.StatusInternalServerError, output)
		return
	}

	// Send the response to the user
	writeJSONToClient(w, http.StatusAccepted, output)
}

func writeJSONToClient(w http.ResponseWriter, status int, output Output) {
	// Send the response to the user
	rawOutput, err := json.Marshal(output)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{}"))
		return
	}
	w.WriteHeader(status)
	_, _ = w.Write(rawOutput)
}
