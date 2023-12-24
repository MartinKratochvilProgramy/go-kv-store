# Creating a custom memory key-value store using Golang

In-memory key-value storage like Redis is an important of most large web applications. The main advantage lies in storing data in memory, making such databases blazing fast, however at the cost of being more expensive to run compared to standard databases, since they do not run using SSD but RAM. Thus, such databases are perfect for storing data that needs to be accessible quickly but does not need to be stored for long periods of time. 

In this article I will describe my approach for creating such system. To implement this idea I have chosen Golang, as it is fast and low-level, while being relatively easy to write effective code. Main features of this project should be:
- Data is highly accessible.
- Stored data has an expiration time, meaning after certain time, it is removed.
- Logs should be written so that after stopping the database, it can be reconstructed from these logs.
- Expose the application through a web server with relevant methods.

## 

Let's first create a project by running:

```
go mod init go-redis
```

Then we can start working on the redis package. Create folder redis

