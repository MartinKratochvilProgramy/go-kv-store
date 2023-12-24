// main.go
package main

import (
	"flag"
	"go-redis/router"
	"go-redis/storage"
	"time"
)

func main() {
	useLogs := flag.Bool("useLogs", true, "If true, logs will be written to logsFilename.")
	reconstructFromLogs := flag.Bool("reconstructFromLogs", false, "Reconstruct database from existing log file.")
	logsFilename := flag.String("logsFilename", "logs.txt", "Log file name.")
	expiration := flag.Duration("expiration", 3*time.Minute, "How long to store key-value pairs for.")
	port := flag.Int("port", 3000, "Port at which to run the server.")
	flag.Parse()

	redis := storage.NewStorage(useLogs, reconstructFromLogs, logsFilename, expiration)

	router := router.NewRouter(port, redis)

	router.Run()

}
