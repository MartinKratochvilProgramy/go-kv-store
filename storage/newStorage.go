package storage

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

type Storage struct {
	mu         sync.Mutex
	expiration time.Duration
	store      *map[string]StoreWrite
	useLogs    bool
	logFile    *os.File
	Tail       *StoreWrite
	Head       *StoreWrite
}

func NewStorage(
	useLogs *bool,
	reconstructFromLogs *bool,
	logsFilename *string,
	expiration *time.Duration,
) *Storage {
	newStore := make(map[string]StoreWrite)

	newLogFile, err := os.OpenFile("./logs/"+*logsFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	newStorage := &Storage{
		expiration: *expiration,
		store:      &newStore,
		useLogs:    *useLogs,
		logFile:    newLogFile,
	}

	if *reconstructFromLogs {
		fmt.Println("Reconstructing from logs...")
		err := newStorage.reconstructFromLogs()
		if err != nil {
			log.Fatal(err)
		}
	}

	ticker := time.Tick(1 * time.Second)
	go func() {
		for {
			<-ticker
			newStorage.cleanupExpiredEntries()
		}
	}()

	return newStorage
}
