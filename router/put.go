package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) put(c *gin.Context) {
	var data map[string]interface{}
	body, err := ioutil.ReadAll(c.Request.Body)

	err = json.Unmarshal(body, &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for key, value := range data {
		err = r.redis.Put(key, value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, "OK")
	return
}
