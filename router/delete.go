package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) delete(c *gin.Context) {
	var body struct {
		Key string `json:"Key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func() {
		err := r.storage.Delete(body.Key)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{body.Key: "Deleted"})
		return
	}()
}
