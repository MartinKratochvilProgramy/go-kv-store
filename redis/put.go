package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

func (redis *Redis) Put(
	key string,
	value interface{},
	id uuid.UUID,
	timestamp time.Time,
) error {
	if time.Since(timestamp) > redis.expiration {
		return nil
	}

	if redis.useLogs {
		valueBytes, err := json.Marshal(value)
		valueString := fmt.Sprintf(
			"PUT, %s, %s, {\"%s\": %s}\n",
			timestamp.Format(time.RFC3339),
			id,
			key,
			string(valueBytes))

		_, err = redis.logFile.WriteString(valueString)
		if err != nil {
			return fmt.Errorf("Failed to write logs: %w", err)
		}
	}

	newStoreWrite := &StoreWrite{
		key:            key,
		id:             id,
		createdAt:      timestamp,
		value:          value,
		prevStoreWrite: redis.Head,
		nextStoreWrite: nil,
	}
	// if new store, set tail
	if redis.Tail == nil {
		redis.Tail = newStoreWrite
	}
	// if new store, set tail
	if redis.Head == nil {
		redis.Head = newStoreWrite
	}

	// set prevStoreWrite next to newStoreWrite
	redis.Head.nextStoreWrite = newStoreWrite

	(*redis.store)[key] = *newStoreWrite
	// set new head
	redis.Head = newStoreWrite

	return nil
}
