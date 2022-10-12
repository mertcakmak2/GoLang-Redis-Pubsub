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

	// redisClient = redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379",
	// })

	// sentinel := redis.NewSentinelClient(&redis.Options{
	// 	Addr:     "34.118.67.177:30160", // 6379:31317/TCP, 26379:30160/TCP
	// 	Password: "4qFSVSH0Y8",
	// })
	// fmt.Println("Masters: ", sentinel.Masters(ctx))

	// addr, err1 := sentinel.GetMasterAddrByName(ctx, "mymaster").Result()
	// if err1 != nil {
	// 	panic("Connection refused " + err1.Error())
	// } else {
	// 	fmt.Println("Master address: ", addr)
	// }

	// rdb := redis.NewFailoverClient(&redis.FailoverOptions{
	// 	MasterName: "mymaster",
	// 	Password:   "4qFSVSH0Y8",
	// 	SentinelAddrs: []string{
	// 		"34.118.67.177:30160",
	// 	},
	// })

	// err1 := rdb.Ping(ctx).Err()
	// if err1 != nil {
	// 	panic("Connection refused " + err1.Error())
	// } else {
	// 	fmt.Println("hata yok")
	// }

	rdb := connectRedis()
	fmt.Println(rdb)
	fmt.Println("get value from redis...")
	val, err := rdb.Get(ctx, "conn-temp").Result()
	if err != nil {
		fmt.Println("no data")
	}
	fmt.Println("Cached value:", val)

	time.Sleep(time.Hour)

	// Kubernetes Redis cluster
	// redisClient = redis.NewClient(&redis.Options{
	// 	Addr: "34.118.67.177:31317",
	// 	// Addr:     addr[0] + ":" + "6379",
	// 	Password: "4qFSVSH0Y8",
	// })

	// fmt.Printf("redisClient.ClientGetName(ctx): %v\n", redisClient.Info(ctx))

	// err := redisClient.Ping(ctx).Err()
	// if err != nil {
	// 	panic("Connection refused " + err.Error())
	// }

	// err = redisClient.Set(ctx, "temp2", "temp value2", time.Second).Err() // time duration or 0
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	panic(err)
	// } else {
	// 	fmt.Println("wrote")
	// }

	// wg := sync.WaitGroup{}
	// wg.Add(20)
	// for i := 0; i < 20; i++ {
	// 	go func(ind int) {
	// 		// val, err := redisClient.Get(ctx, "temp").Result()
	// 		// if err != nil {
	// 		// 	panic(err)
	// 		// }
	// 		// fmt.Println("Cached value:", val)
	// 		setValue("1")
	// 		setValue("2")
	// 		setValue("3")
	// 		// setValue("temp2")
	// 		wg.Done()
	// 	}(i)
	// }
	// wg.Wait()

	// val, err := redisClient.Get(ctx, "2").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Cached value:", val)

	// var channelName = "test-channel"
	// go producer(channelName)
	// go consumer(channelName)

	// time.Sleep(time.Hour * 1)
}

func connectRedis() *redis.Client {
	fmt.Println("Attempt to connection...")
	redisClient = redis.NewClient(&redis.Options{
		Addr: "34.118.67.177:31317",
		// Addr:     addr[0] + ":" + "6379",
		Password: "4qFSVSH0Y8",
	})

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		panic("Connection refused " + err.Error())
	}

	err = redisClient.Set(ctx, "conn-temp", "conn-temp", time.Second).Err() // time duration or 0
	if err != nil {
		fmt.Println("Error: ", err)
		return connectRedis() // Blocking for connect to slave node
	}
	return redisClient
}

func setValue(key string) {
	err := redisClient.Set(ctx, key, time.Now().Format("15:04:05"), time.Minute*1).Err() // time duration or 0
	if err != nil {
		panic(err)
	}
	fmt.Println("wrote")
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
