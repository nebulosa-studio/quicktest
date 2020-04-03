package main

import (
	"fmt"
	"nebulosa-studio/quicktest/status"

	"github.com/go-redis/redis"
)

func main() {
	redisTest("localhost", "6379")
}

func redisTest(host, port string) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
	r, err := client.Ping().Result()
	if err != nil {
		status.Print("redis", status.Error, err.Error())
		return
	}

	status.Print("redis", status.Success, r)
}
