package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

var errSaveNewSubscriber = errors.New("can't save new subscriber, try again later")
var errInvalidMembershipType = errors.New("not a valid membership type")

var membershipTypes = []string{"trial", "standard", "premium"}

type Input struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	MembershipType string `json:"membership_type"`
	UID            string `json:"uid"`
}

type Output struct {
	SubscribedAt time.Time `json:"subscribed_at"`
	UID          string    `json:"uid"`
	Errors       []string  `json:"errors"`
}

func addNewSubscriber(db *sql.DB, rdb *redis.Client, input Input) (Output, error) {
	// Validate the input
	if !validMembership(input.MembershipType) {
		return Output{
			UID:    input.UID,
			Errors: []string{errInvalidMembershipType.Error()},
		}, errInvalidMembershipType
	}

	err := saveNewSubscriber(db, input.Name, input.Email, input.MembershipType, input.UID, time.Now())
	if err != nil {
		return Output{
			UID:    input.UID,
			Errors: []string{errSaveNewSubscriber.Error()},
		}, err
	}

	cmd := rdb.RPush(context.Background(), "subscribers", input.Email)
	if cmd.Err() != nil {
		return Output{}, cmd.Err()
	}
	return Output{
		UID:          input.UID,
		SubscribedAt: time.Now(),
	}, nil
}

func validMembership(membership string) bool {
	for _, m := range membershipTypes {
		if m == membership {
			return true
		}
	}
	return false
}
