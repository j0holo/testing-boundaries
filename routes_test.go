package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-txdb"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/matryer/is"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	txdb.Register("txdb", "postgres", "host=localhost port=5432 user=user password=password dbname=my_database sslmode=disable")

	os.Exit(m.Run())
}

func getDatabaseConnection(t *testing.T, createTables bool) *sql.DB {
	t.Helper()

	// Hardcode the dsn for this simple example
	db, err := sql.Open("txdb", "host=localhost port=5432 user=user password=password dbname=my_database sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	if createTables {
		err = createSubscribersTable(db)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Run this function after the test has ended that called this function.
	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func getRedisDatabaseConnection(t *testing.T) *redis.Client {
	t.Helper()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		rdb.Close()
	})

	return rdb
}

func startTestServer(t *testing.T, createTables bool) *httptest.Server {
	t.Helper()
	app := &application{db: getDatabaseConnection(t, createTables), rdb: getRedisDatabaseConnection(t)}
	return httptest.NewServer(newRouter(app))
}

func getJsonOutput(body io.ReadCloser) (*Output, error) {
	output := Output{}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &output)
	return &output, err
}

func TestJsonHandlerWithNoInput(t *testing.T) {
	is := is.New(t)
	server := startTestServer(t, true)
	defer server.Close()

	response, err := http.Post(server.URL+"/api/v1/json", "application/json", nil)
	is.NoErr(err)

	is.Equal(response.StatusCode, http.StatusBadRequest)
}

func TestJsonHandler(t *testing.T) {
	is := is.New(t)
	server := startTestServer(t, true)
	defer server.Close()

	input := Input{
		UID:            "some-random-uid",
		Email:          "john@example.com",
		Name:           "john",
		MembershipType: "trial",
	}

	data, err := json.Marshal(input)
	is.NoErr(err)

	response, err := http.Post(server.URL+"/api/v1/json", "application/json", bytes.NewReader(data))
	is.NoErr(err)

	output, err := getJsonOutput(response.Body)
	is.NoErr(err)

	expect := &Output{
		SubscribedAt: output.SubscribedAt,
		UID:          input.UID,
		Errors:       nil,
	}

	is.Equal(response.StatusCode, http.StatusAccepted)
	is.Equal(output, expect)
}

func TestJsonHandlerWitNoSchema(t *testing.T) {
	is := is.New(t)
	server := startTestServer(t, false)
	defer server.Close()

	input := Input{
		UID:            "some-random-uid",
		Email:          "john@example.com",
		Name:           "john",
		MembershipType: "trial",
	}

	data, err := json.Marshal(input)
	is.NoErr(err)

	response, err := http.Post(server.URL+"/api/v1/json", "application/json", bytes.NewReader(data))
	is.NoErr(err)

	output, err := getJsonOutput(response.Body)
	is.NoErr(err)

	expect := &Output{
		SubscribedAt: output.SubscribedAt,
		UID:          input.UID,
		Errors:       []string{errSaveNewSubscriber.Error()},
	}

	is.Equal(response.StatusCode, http.StatusInternalServerError)
	is.Equal(output, expect)
}

func TestJsonHandlerWithWrongSubscription(t *testing.T) {
	is := is.New(t)
	server := startTestServer(t, true)
	defer server.Close()

	input := Input{
		UID:            "some-random-uid",
		Email:          "john@example.com",
		Name:           "john",
		MembershipType: "this-does-not-exist",
	}

	data, err := json.Marshal(input)
	is.NoErr(err)

	response, err := http.Post(server.URL+"/api/v1/json", "application/json", bytes.NewReader(data))
	is.NoErr(err)

	output, err := getJsonOutput(response.Body)
	is.NoErr(err)

	expect := &Output{
		UID:    input.UID,
		Errors: []string{errInvalidMembershipType.Error()},
	}

	is.Equal(response.StatusCode, http.StatusBadRequest)
	is.Equal(output, expect)
}
