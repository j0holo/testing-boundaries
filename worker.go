package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
)

func readEmails(rdb *redis.Client, queue string, w io.Writer) {
	result, err := rdb.LPop(context.Background(), queue).Result()
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(result))
	if err != nil {
		log.Println(err)
	}
}
