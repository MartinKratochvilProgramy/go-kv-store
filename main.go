// main.go
package main

import (
	"bufio"
	"fmt"
	"go-redis/redis"
	"go-redis/router"
	"os"
	"strings"
)

func main() {
	redis := redis.NewRedis()
	router := router.NewRouter("3000", redis)

	router.Run()

	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Sprintf("Error reading input: %v", err))
		}

		input = strings.TrimSpace(input)
		commands := strings.Fields(input)

		command := commands[0]
		key := commands[1]
		value := strings.Join(commands[2:], " ")

		if input == "q" {
			fmt.Println("Exiting the program.")
			break
		}

		switch command {
		case "GET":
			value := redis.Get(key)
			fmt.Println(value)
		case "PUT":
			redis.Put(key, value)
		default:
			fmt.Println("Wrong command! Try 'GET key' or 'PUT key value'.")
		}
	}
}
