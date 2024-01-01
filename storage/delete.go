package storage

import (
	"errors"
	"fmt"
)

func (storage *Storage) Delete(
	key string,
) error {

	storage.mu.Lock()

	if storage.useLogs {
		valueString := fmt.Sprintf(
			"DELETE, %s\n",
			key,
		)

		_, err := storage.logFile.WriteString(valueString)
		if err != nil {
			return fmt.Errorf("Failed to write logs: %w", err)
		}
	}

	storeWrite, ok := (*storage.store)[key]
	if !ok {
		storage.mu.Unlock()
		return errors.New("Key not found")
	}

	if storage.Tail.key == key && storage.Tail.nextStoreWrite != &storeWrite {
		storage.Tail = storeWrite.nextStoreWrite
	}

	if storage.Head.key == key && storage.Head.nextStoreWrite != &storeWrite {
		storage.Head = storeWrite.prevStoreWrite
	}

	if storeWrite.prevStoreWrite != nil && storeWrite.nextStoreWrite != nil {
		storeWrite.prevStoreWrite.nextStoreWrite = storeWrite.nextStoreWrite
	}

	delete(*storage.store, key)

	storage.mu.Unlock()
	return nil
}
