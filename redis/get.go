package redis

func (redis *Redis) Get(key string) interface{} {
	storeWrite, ok := (*redis.store)[key]

	if !ok {
		return nil
	}

	return storeWrite.value
}
