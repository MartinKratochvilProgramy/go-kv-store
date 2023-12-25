package storage

import (
	"time"
)

func (storage *Storage) cleanupExpiredEntries() {
	var currentStoreWrite = storage.Tail
	if currentStoreWrite == nil {
		return
	}

	storage.mu.Lock()
	for {
		if currentStoreWrite == nil {
			storage.mu.Unlock()
			return
		}
		if time.Since(currentStoreWrite.createdAt) > storage.expiration {
			// if write is expired, delete it
			storage.Tail = currentStoreWrite.nextStoreWrite
			currentStoreWrite.prevStoreWrite = nil

			delete(*storage.store, currentStoreWrite.key)

			currentStoreWrite = currentStoreWrite.nextStoreWrite
		} else {
			// if write is not expired, break the loop
			storage.mu.Unlock()
			return
		}
	}
}
