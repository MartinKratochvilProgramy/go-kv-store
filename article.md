# Creating a custom memory key-value store using Golang

In-memory key-value storage is an important part of most large web applications. The main advantage lies in storing data in memory, making such databases ✨blazing fast✨, however at the cost of being more expensive to run compared to standard databases since they do not run using SSD but RAM. Thus, such databases are perfect for storing data that needs to be accessible quickly, but does not need to be stored for long periods of time. 

In this article I will describe my approach for creating such system. To implement this idea I have chosen Golang, as it is relatively fast and low-level enough, while being easy to write effective code in. Main features of this project are:

- Data should be quickly accessible - achieve sub-ms speed for data lookup.
- Stored data has an expiration time, meaning after certain duration, it is automatically removed.
- Expose the application through a web server with relevant methods.

## 

Let's first create a project by running:

```
go mod init go-kv-store
```

Then we can start working on the storage. Create folder ./storage, inside we will make a Storage struct with store map - the reason for using map is that it allows us quickly access unstructured data by given key. Map keys will be strings and values will store in an empty interface StoreWrite - each write will have a unique identifier, a timestamp so we can delete the expired values and the actual data. We also have to include mutex in the storage struct - since our server can handle concurrent requests, we will have to lock the storage before every read and write so as to not get any duplicate write. The stored data will be defined as an empty interface so that we are not limited by data types.

``` go
// storage/newStorage.go

package storage

import (
	"time"
	"sync"
	"github.com/gofrs/uuid"
)

type StoreWrite struct {
	id        uuid.UUID
	createdAt time.Time
	value     interface{}
}

type Storage struct {
	mu         sync.Mutex
	store      *map[string]StoreWrite
}

func NewStorage() *Storage {
	newStore := make(map[string]StoreWrite)

	newStorage := &Storage{
		store: &newStore,
	}

	return newStorage
}

```

Using map as opposed to other data structures has one big advantage - it allows us to access unstructured data fairly quicky if we know the key that needs to be looked up, no need for loops. We then create a simple get method that return value by given key, but everytime we get an item, we need to lock the mutex:
``` go
// storage/get.go

package storage

func (storage *Storage) Get(key string) interface{} {
	storage.mu.Lock()
	storeWrite, ok := (*storage.store)[key]
	storage.mu.Unlock()

	if !ok {
		return nil
	}

	return storeWrite.value
}
```

Just get will not be enough, so create a put method that allows us insert value into our map - again we have to lock the mutex before inserting:
```go
// storage/put.go

package storage

import (
	"time"

	"github.com/gofrs/uuid"
)

func (storage *Storage) Put(
	key string,
	value interface{},
	id uuid.UUID,
	timestamp time.Time,
) error {

	newStoreWrite := &StoreWrite{
		id:        id,
		createdAt: timestamp,
		value:     value,
	}

	storage.mu.Lock()
	(*storage.store)[key] = *newStoreWrite
	storage.mu.Unlock()

	return nil
}
```

Simple enough. Now that we have basic storage functionality, we can start working on how to expose the API. Although Go is not an OOP language by design, I find that the best way to structure a REST API is by creating a new struct and passing necessary code as a dependency injection. Go has built-in http package that works well out of the box, in this case however, I will be using the [Gin Web Framework](https://github.com/gin-gonic/gin). Create a router package with NewRouter method:

``` go
// router/newRouter.go

package router

import (
	"fmt"
	"go-kv-store/storage"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	storage    *storage.Storage
	router     *gin.Engine
	serverAddr *string
}

func NewRouter(port *int, storage *storage.Storage) *Router {
	serverAddr := "127.0.0.1:" + fmt.Sprint(*port)

	router := gin.Default()

	return &Router{
		storage:    storage,
		router:     router,
		serverAddr: &serverAddr,
	}
}
```

Now we can create handlers for get and put methods. Since we passed the storage as dependency, it can be easily accessed from the router - as opposed to awkwardly passing it around. We put the sotrage.Get method inside a goroutine, so that the server can handle concurrent read and writes without any conflicts:

``` go
// router/get.go

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) get(c *gin.Context) {
	var body struct {
		Key string `json:"Key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func() {
		value := r.storage.Get(body.Key)

		if value == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found."})
			return
		}

		c.JSON(http.StatusOK, gin.H{body.Key: value})
		return
	}()
}
```

``` go
// router/put.go

package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func (r *Router) put(c *gin.Context) {
	var data map[string]interface{}
	body, err := ioutil.ReadAll(c.Request.Body)

	err = json.Unmarshal(body, &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for key, value := range data {
		id, err := uuid.NewV4()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		go func(key string, value interface{}) {
			err = r.storage.Put(key, value, id, time.Now())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}(key, value)
	}

	c.JSON(http.StatusCreated, "OK")
	return

}
```

Lastly create run method where each route will be registered and we can run the API:

``` go
// router/run.go

package router

func (r *Router) Run() {
	r.router.GET("/", r.get)
	r.router.PUT("/", r.put)

	r.router.Run(*r.serverAddr)
}
```

Now that all the basic functionality is defined, we can start the application by piecing all the code together in main function, we also add a parser, so that port can be specified if necessary:
``` go
// main.go
package main

import (
	"flag"
	"go-kv-store/router"
	"go-kv-store/storage"
	"time"
)

func main() {
	port := flag.Int("port", 3000, "Port at which to run the server.")
	flag.Parse()

	redis := storage.NewStorage()

	router := router.NewRouter(port, redis)
	router.Run()
}
```

We can run the code by running:

```shell
go run . --port=3000
```

Now we can put data into the storage with a PUT request. Since we are using an empty interface, the storage can handle even more complex data types:
```http
PUT http://127.0.0.1:3000/

{
  "foo": "bar",
  "complexVal": {"foo": 0, "bar": 3.14}
}
```
Calling
```http
GET http://127.0.0.1:3000/

{
  "key": "complexVal"
}
```

Yields:
```json
{
  "complexVal": {
    "bar": 3.14,
    "foo": 0
  }
}
```

Great! We can store data. But how about removing the data after it is expired?

## Linked list to the rescue

We would like the database writes to expire after some time - easy enough, the timestamp is already defined. But here's the beef: how do we identify the values for deletion? Definetly some kind of loop will be necessary, but looping through all the values is terribly inefficent. We can assume that data written into the database will be sorted chronologically, so maybe putting the data into a list would make sense, since we can just slice the list at given value where timestamp is not expired anymore and not have to iterate further, but using list would be slow for lookup - the main reason I chose map is because there is no need to iterate over when searching for values.

Another solution could be to keep the data in the map and list at the same time, where map would be used for lookup and list to find data older than some value, but duplicating data in a database is a terrible idea.

The best way I could think of to store time-structured data to a map is to implement a linked list data structure where each item contains not only the necessary data, but also a pointer to writes before and after it. Again, I am assuming that data will be written to the database in a chronological order and the order will not change (ie. the timestamp is final and determined by the application).

Let's tweak the storage struct and add necessary pointers. We also need to keep track of head and tail, so that we know where to start when traversing the list:
``` go
// storage/newStorage.go

...

type StoreWrite struct {
	key            string
	id             uuid.UUID
	createdAt      time.Time
	value          interface{}
	prevStoreWrite *StoreWrite
	nextStoreWrite *StoreWrite
}

type Storage struct {
	mu         sync.Mutex
	expiration time.Duration
	store      *map[string]StoreWrite
	useLogs    bool
	logFile    *os.File
	Tail       *StoreWrite
	Head       *StoreWrite
}

func NewStorage(expiration *time.Duration,) *Storage { 
    ...

    newStorage := &Storage{
		expiration: *expiration,
        ....
    }

    ...
}
```

We can pass expiration time as an env variable from the parser:

``` go
// main.go

func main() {
    expiration := flag.Duration("expiration", 3*time.Minute, "How long to store key-value pairs for.")

    ...

	redis := storage.NewStorage(expiration)

    ...
}
```

Modify the put method in the storage struct so that each new write has prevStoreWrite equal to storage.Head and move the head:

``` go
// storage/put.go

package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

func (storage *Storage) Put(
	key string,
	value interface{},
	id uuid.UUID,
	timestamp time.Time,
) error {
	if time.Since(timestamp) > storage.expiration {
		return nil
	}

	newStoreWrite := &StoreWrite{
		key:            key,
		...
        prevStoreWrite: storage.Head,
		nextStoreWrite: nil,
	}
	// if empty store, set tail
	if storage.Tail == nil {
		storage.Tail = newStoreWrite
	}

	if storage.Head == nil {
		// if empty store, set head
		storage.Head = newStoreWrite
	} else {
		// set new head only if it exists
		storage.Head.nextStoreWrite = newStoreWrite
	}

	// set new head
	storage.Head = newStoreWrite

	(*storage.store)[key] = *newStoreWrite

    return nil
}
```

Now we can traverse the linked list, look for expired values and when an unexpired value is encountered, break the loop:

``` go
// storage/cleanupExpiredEntries.go

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
```

Start the cleanup function in a goroutine so that it does not interfere with the storage functionality:

```go
// storage/newStorage.go
func NewStorage(expiration *time.Duration,) *Storage { 
    ...

    ticker := time.Tick(1 * time.Second)
	go func() {
		for {
			<-ticker
			newStorage.cleanupExpiredEntries()
		}
	}()

    ...
}
```

Lastly we will pass the expiration parameter as an environment variable in the main function:

```go
// main.go

func main() {
	...
	expiration := flag.Duration("expiration", 3*time.Minute, "How long to store key-value pairs for.")

	redis := storage.NewStorage(, expiration)
	...
}

```

## Deleting

Although we have already designed a way to automatically delete writes, let's also implement a way for the user to delete - what kind of database would it be if there was no option to delete? In the storage struct create a delete method - again we have to lock the mutex to avoid conflicts. We also have to restructure the linked list on deletion. First check if the deleted key is either head or tail and if it is, move the nextStoreWrite or prevStoreWrite, then connect the previous and next writes by moving the relevant pointers.

```go
package storage

import (
	"errors"
)

func (storage *Storage) Delete(key string) error {

	storage.mu.Lock()

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
```

Another problem we run into is when editing existing write, we do not edit the timestamp, the app does it for us. Consequently, the edited write, that is changed by the Put method stays at the same place in the linked list, but we would like to move it to the end, so it would not expire immediatedly. We should first check if a write exsits and if it does, we can use the newly defined Delete method to remove it from the list, and then append new one to the end. Alter the Put method like this:

```go
// storage/put.go

func (storage *Storage) Put(
	key string,
	value interface{},
	id uuid.UUID,
	timestamp time.Time,
) error {
	if time.Since(timestamp) > storage.expiration {
		return nil
	}

	storage.mu.Lock()
	_, ok := (*storage.store)[key]
	storage.mu.Unlock()
	if ok {
		// if key exists, remove it and append to the end
		// this is to update it's timestamp
		storage.Delete(key)
	}

	storage.mu.Lock()

	...

}
```

## Stress test

Let's test how much load can our storage handle. Run the sotrage and this time set the expiraton variable to 1s - this way we are guaranteed that some entries will expire during the test so that there has to be some work done by the cleanup routine: 

```shell
go run . --port=3000 --expiration=1s
```

I have also written a python script to simulate the server load - using ProcessPoolExecutor we spawn 8 concurrent processes which send 100 000 put and get requests, essentially inserting and reading every key - we then time the total execution time for all these 200 000 requests. The script is provided below:

```python
import time
import requests
import json
from concurrent.futures import ProcessPoolExecutor
from urllib3.util.retry import Retry
from requests.adapters import HTTPAdapter

session = requests.Session()
retry = Retry(connect=3, backoff_factor=0.5)
adapter = HTTPAdapter(max_retries=retry)
session.mount('http://', adapter)

def send_put_request(start, end):
    url = "http://127.0.0.1:3000/"
    for i in range(start, end):
        data = {str(i): str(i)}
        _ = session.put(url, data=json.dumps(data))

def send_get_request(start, end):
    url = "http://127.0.0.1:3000/"
    for i in range(start, end):
        data = {"key": str(i)}
        _ = session.get(url, data=json.dumps(data))

def main():
    N = 100_000
    processes = 8
    step = N // processes

    with ProcessPoolExecutor(max_workers=processes) as executor:
        ranges = [(i * step, (i + 1) * step) for i in range(processes)]

        put_futures = [executor.submit(send_put_request, start, end) for start, end in ranges]
        get_futures = [executor.submit(send_get_request, start, end) for start, end in ranges]

        ts = time.time()    
        for future in put_futures + get_futures:
            future.result()
        te = time.time()
        print(f"Exec time: {te-ts}s")

if __name__ == "__main__":
    main()
```

Here's the result: on my machine (personal laptop with Intel i7, the process consumed around 700MB of RAM) total execution time was 58.283s, which means 0.296ms per request. We are well under the 1ms limit, considering the amount of work it took, that's pretty good!

## Conclusion

In this article we implemented a simple key-value store in Go. It turns out that Go is a great choice if we want to write performant and relatively low-level code easily. Our implemented key value store can handle thousands of requests while keeping the response time unde 1ms.

Of course this project cannot fully replace existing solutions such as Redis, this has been mainly an excercise to implement such solution and the corresponding data structures.

You can find the complete code on my [GitHub](https://github.com/MartinKratochvilProgramy/go-kv-store).