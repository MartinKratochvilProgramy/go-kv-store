package redis

import (
	"time"
)

func (redis *Redis) cleanupExpiredEntries() {
	var currentStoreWrite = redis.Tail
	if currentStoreWrite == nil {
		return
	}

	redis.mu.Lock()
	for {
		if currentStoreWrite == nil {
			redis.mu.Unlock()
			return
		}
		if time.Since(currentStoreWrite.createdAt) > redis.expiration {
			redis.Tail = currentStoreWrite.nextStoreWrite
			currentStoreWrite.prevStoreWrite = nil
			delete(*redis.store, currentStoreWrite.key)
		}
		currentStoreWrite = currentStoreWrite.nextStoreWrite
	}
}
