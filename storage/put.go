package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

func (storage *Storage) Put(
	key string,
	value interface{},
	id uuid.UUID,
	timestamp time.Time,
) error {
	if time.Since(timestamp) > storage.expiration {
		return nil
	}

	if storage.useLogs {
		valueBytes, err := json.Marshal(value)
		valueString := fmt.Sprintf(
			"PUT, %s, %s, {\"%s\": %s}\n",
			timestamp.Format(time.RFC3339),
			id,
			key,
			string(valueBytes))

		_, err = storage.logFile.WriteString(valueString)
		if err != nil {
			return fmt.Errorf("Failed to write logs: %w", err)
		}
	}

	newStoreWrite := &StoreWrite{
		key:            key,
		id:             id,
		createdAt:      timestamp,
		value:          value,
		prevStoreWrite: storage.Head,
		nextStoreWrite: nil,
	}
	// if new store, set tail
	if storage.Tail == nil {
		storage.Tail = newStoreWrite
	}
	// if new store, set tail
	if storage.Head == nil {
		storage.Head = newStoreWrite
	}

	// set prevStoreWrite next to newStoreWrite
	storage.Head.nextStoreWrite = newStoreWrite

	(*storage.store)[key] = *newStoreWrite
	// set new head
	storage.Head = newStoreWrite

	return nil
}
