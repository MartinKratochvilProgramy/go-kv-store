package storage

import "fmt"

func (storage *Storage) GetAll() []interface{} {
	fmt.Println("Get all")
	var currentStoreWrite = storage.Tail
	// if currentStoreWrite == nil {
	// 	return nil
	// }

	for {
		if currentStoreWrite == nil {
			return nil
		}

		fmt.Println("value: ", currentStoreWrite.value, " key: ", currentStoreWrite.key, " timestamp: ", currentStoreWrite.createdAt)
		currentStoreWrite = currentStoreWrite.nextStoreWrite
	}
}
