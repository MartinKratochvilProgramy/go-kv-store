package storage

func (storage *Storage) Get(key string) interface{} {
	storeWrite, ok := (*storage.store)[key]

	if !ok {
		return nil
	}

	return storeWrite.value
}
