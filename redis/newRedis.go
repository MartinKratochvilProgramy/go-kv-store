package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

type StoreWrite struct {
	id        uuid.UUID
	createdAt time.Time
	value     interface{}
}

type Redis struct {
	store   *map[string]StoreWrite
	logFile *os.File
}

func NewRedis() *Redis {
	newStore := make(map[string]StoreWrite)

	newLogFile, err := os.OpenFile("./logs/logs.txt", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
	}

	return &Redis{
		store:   &newStore,
		logFile: newLogFile,
	}
}
