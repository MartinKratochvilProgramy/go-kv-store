package router

import (
	"fmt"
	"go-redis/redis"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	redis      *redis.Redis
	router     *gin.Engine
	serverAddr *string
}

func NewRouter(port string, redis *redis.Redis) *Router {
	serverAddr := fmt.Sprint("127.0.0.1:", port)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Router{
		redis:      redis,
		router:     router,
		serverAddr: &serverAddr,
	}
}
