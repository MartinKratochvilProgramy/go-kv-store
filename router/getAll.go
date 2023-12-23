package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) getAll(c *gin.Context) {
	r.redis.GetAll()

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
	return
}
