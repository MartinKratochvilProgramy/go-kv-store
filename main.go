// main.go
package main

import (
	"bufio"
	"fmt"
	"go-redis/redis"
	"os"
	"strings"
)

func main() {
	redis := redis.NewRedis()

	fmt.Println("Listening...")

	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Sprintf("Error reading input: %v", err))
		}

		input = strings.TrimSpace(input)
		commands := strings.Fields(input)

		if input == "q" {
			fmt.Println("Exiting the program.")
			break
		}

		switch commands[0] {
		case "GET":
			value := redis.Get(commands[1])
			fmt.Println(value)
		case "PUT":
			redis.Put(commands[1], commands[2])
		default:
			fmt.Println("Wrong command! Try 'GET key' or 'PUT key value'.")
		}
	}
}
