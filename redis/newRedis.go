package redis

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

type StoreWrite struct {
	key            string
	id             uuid.UUID
	createdAt      time.Time
	value          interface{}
	prevStoreWrite *StoreWrite
	nextStoreWrite *StoreWrite
}

type Redis struct {
	mu         sync.Mutex
	expiration time.Duration
	store      *map[string]StoreWrite
	logFile    *os.File
	Tail       *StoreWrite
	Head       *StoreWrite
}

func NewRedis() *Redis {
	newStore := make(map[string]StoreWrite)

	newLogFile, err := os.OpenFile("./logs/logs.txt", os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}

	newRedis := &Redis{
		expiration: 25 * time.Second,
		store:      &newStore,
		logFile:    newLogFile,
	}

	newRedis.reconstructFromLogs()

	ticker := time.Tick(1 * time.Second)
	go func() {
		for {
			<-ticker
			// newRedis.GetAll()
			newRedis.cleanupExpiredEntries()
		}
	}()

	return newRedis
}
