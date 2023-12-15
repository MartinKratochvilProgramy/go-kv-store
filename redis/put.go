package redis

import (
	"time"

	"github.com/gofrs/uuid"
)

func (redis *Redis) Put(key string, value interface{}) error {

	newId, err := uuid.NewV4()
	if err != nil {
		return err
	}

	newStoreWrite := &StoreWrite{
		id:        newId,
		createdAt: time.Now(),
		value:     value,
	}

	(*redis.store)[key] = *newStoreWrite

	return nil
}
