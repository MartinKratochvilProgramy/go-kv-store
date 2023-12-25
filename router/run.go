package router

import (
	"fmt"
)

func (r *Router) Run() {
	// r.router.GET("/", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, "OK")
	// })

	r.router.GET("/", r.get)
	// r.router.GET("/get-all", r.getAll)
	r.router.PUT("/", r.put)

	fmt.Println("Server running on", *r.serverAddr)
	r.router.Run(*r.serverAddr)
}
