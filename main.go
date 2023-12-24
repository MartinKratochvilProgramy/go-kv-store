// main.go
package main

import (
	"go-redis/redis"
	"go-redis/router"
)

func main() {
	redis := redis.NewRedis()
	// redis.GetAll()

	router := router.NewRouter("3000", redis)

	router.Run()

}
