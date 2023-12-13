package redis

import (
	"time"

	"github.com/gofrs/uuid"
)

type StoreWrite struct {
	id        uuid.UUID
	createdAt time.Time
	data      interface{}
}

type Redis struct {
	store *map[string]StoreWrite
}

func NewRedis() *Redis {
	newStore := make(map[string]StoreWrite)
	return &Redis{
		store: &newStore,
	}
}
