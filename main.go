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
	go producer(channelName)
	go consumer(channelName)

	time.Sleep(time.Hour * 1)
}

func producer(channelName string) {
	for range time.Tick(time.Second * 3) {
		time := time.Now().Format("15:04:05")
		fmt.Println("Sent message: " + time)
		redisClient.Publish(ctx, channelName, time)
	}
}

func consumer(channelName string) {
	subs := redisClient.Subscribe(ctx, channelName)

	for msg := range subs.Channel() {
		fmt.Println("Received message: ", msg.Payload)
		fmt.Println("Channel Name: ", msg.Channel)
	}
}
