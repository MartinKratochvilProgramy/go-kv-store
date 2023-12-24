package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) getAll(c *gin.Context) {
	r.storage.GetAll()

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
	return
}
