# About this project

This is a simple implementation of key-value store in Go. It serves as simple caching server that can store data for specified duration with sub 1ms read and write speed.

Detailed description of the code and design decisions I took when writing it is in this [Medium article](https://medium.com/@martvil96/creating-an-in-memory-key-value-store-using-golang-2db73aaec087).

# Installation

Clone the repository and run:

```shell
go run .
```

# Build

You can create optimized build by running:

```shell
go build .
```

# How to use

After building and running the application, the API will be exposed on port 3000, the storage has 3 methods, PUT, GET, DELETE

- PUT - creates new write, accepts json with key-value pairs, value can be any data compatible with json:
```http
PUT http://127.0.0.1:3000/

{
  "foo": "bar",
  "complexVal": {"foo": 0, "bar": 3.14}
}
```
- GET - reads data from storage by key
```http
GET http://127.0.0.1:3000/

{
  "key": "complexVal"
}
```
returns:
```json
{
  "complexVal": {
    "bar": 3.14,
    "foo": 0
  }
}
```

- DELETE - deletes wrire:
```http
DELETE http://127.0.0.1:3000/

{
  "key": "complexVal"
}
```

# Expiration

Main usage of the store is intended to be caching, therefore the data will expire after specified time (default 3min).

# Logs

During operation the app can write logs into a .txt file which then can be read to reconstruct the database if it were to go offline.

# Environmental variables

The applications includes number of variables that can be set by the user. The variables are:

- use_logs - boolean - ff true, logs will be written to logs_filename
- reconstruct_from_logs - boolean - reconstruct database from existing log file
- logs_filename - string - name of the log file, has to be .txt
- expiration - time.Duration - how long to store key-value pairs for
- port - int - port at which to run the server


Example run:
```shell
go run . --use_logs=false --reconstruct_from_logs=false --logs_filename=logs.txt --expiration=1m30s --port=3000
```
