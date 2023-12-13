package redis

func (redis *Redis) Get(key string) interface{} {
	value, ok := (*redis.store)[key]

	if !ok {
		return nil
	}

	return value.data
}
