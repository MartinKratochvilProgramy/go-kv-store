package redis

import "fmt"

func (redis *Redis) GetAll() []interface{} {
	fmt.Println("Get all")
	var currentStoreWrite = redis.Tail
	// if currentStoreWrite == nil {
	// 	return nil
	// }

	for {
		if currentStoreWrite == nil {
			return nil
		}

		fmt.Println("value: ", currentStoreWrite.value, " key: ", currentStoreWrite.key)
		currentStoreWrite = currentStoreWrite.nextStoreWrite
	}
}
