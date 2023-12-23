package redis

import (
	"fmt"
)

func (redis *Redis) GetAll() []interface{} {
	for key, value := range *redis.store {
		fmt.Printf("Key: %s, Value: %+v\n", key, value.createdAt)
	}

	return nil
}
