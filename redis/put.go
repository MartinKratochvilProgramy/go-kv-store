package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

func (redis *Redis) Put(key string, value string) error {
	var data interface{}
	err := json.Unmarshal([]byte(value), &data)
	fmt.Println(value, data, err)
	if err != nil {
		return err
	}

	newId, err := uuid.NewV4()
	if err != nil {
		return err
	}

	newStoreWrite := &StoreWrite{
		id:        newId,
		createdAt: time.Now(),
		data:      data,
	}

	(*redis.store)[key] = *newStoreWrite

	return nil
}
