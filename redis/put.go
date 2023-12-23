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

	newStoreWrite := &StoreWrite{
		id:        id,
		createdAt: timestamp,
		value:     value,
	}

	(*redis.store)[key] = *newStoreWrite

	return nil
}
