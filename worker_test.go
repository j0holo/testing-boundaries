package main

import (
	"context"
	"github.com/matryer/is"
	"testing"
)

var queue = "queue_test"

func TestReadEmails(t *testing.T) {
	is := is.New(t)
	rdb := getRedisDatabaseConnection(t)

	testEmail := "test@example.com"

	err := rdb.RPush(context.Background(), queue, testEmail).Err()
	is.NoErr(err)

	wm := &writerMock{emails: []string{}}
	readEmails(rdb, queue, wm)

	is.True(len(wm.emails) == 1)
	is.Equal(wm.emails[0], testEmail)
}

func TestReadEmailsEmptyQueue(t *testing.T) {
	is := is.New(t)
	rdb := getRedisDatabaseConnection(t)

	wm := &writerMock{emails: []string{}}
	readEmails(rdb, queue, wm)

	is.True(len(wm.emails) == 0)
}

type writerMock struct {
	emails []string
}

func (w *writerMock) Write(p []byte) (int, error) {
	w.emails = append(w.emails, string(p))
	return len(p), nil
}
