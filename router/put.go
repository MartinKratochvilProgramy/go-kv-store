package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) put(c *gin.Context) {
	var body struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.redis.Put(body.Key, body.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key not found."})
		return
	}

	c.JSON(http.StatusCreated, "OK")
	return
}
