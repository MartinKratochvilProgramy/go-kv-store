package redis

import (
	"fmt"
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
	useLogs    bool
	logFile    *os.File
	Tail       *StoreWrite
	Head       *StoreWrite
}

func NewRedis(
	useLogs *bool,
	reconstructFromLogs *bool,
	logsFilename *string,
	expiration *time.Duration,
) *Redis {
	newStore := make(map[string]StoreWrite)

	newLogFile, err := os.OpenFile("./logs/"+*logsFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	newRedis := &Redis{
		expiration: *expiration,
		store:      &newStore,
		useLogs:    *useLogs,
		logFile:    newLogFile,
	}

	if *reconstructFromLogs {
		fmt.Println("Reconstructing from logs...")
		err := newRedis.reconstructFromLogs()
		if err != nil {
			log.Fatal(err)
		}
	}

	ticker := time.Tick(1 * time.Second)
	go func() {
		for {
			<-ticker
			newRedis.cleanupExpiredEntries()
		}
	}()

	return newRedis
}
