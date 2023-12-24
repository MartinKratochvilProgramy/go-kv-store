package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) Run() {
	r.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})

	r.router.GET("/get", r.get)
	r.router.GET("/get-all", r.getAll)
	r.router.PUT("/put", r.put)

	r.router.Run(*r.serverAddr)
}
