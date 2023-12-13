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

	value := r.redis.Get(body.Key)
	if value == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found."})
		return
	}

	c.JSON(http.StatusOK, gin.H{body.Key: value})
	return
}
