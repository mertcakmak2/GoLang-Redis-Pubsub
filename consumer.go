package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

func main() {

	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		panic("Connection refused " + err.Error())
	}

	var channelName = "test-channel"
	go consumer(channelName)

	time.Sleep(time.Hour * 1)
}

func consumer(channelName string) {
	subs := redisClient.Subscribe(ctx, channelName)

	for msg := range subs.Channel() {
		fmt.Println("Received message on consumer.go: ", msg.Payload)
		fmt.Println("Channel Name on consumer.go: ", msg.Channel)
	}
}
