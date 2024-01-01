package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
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
		id, err := uuid.NewV4()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		go func(key string, value interface{}) {
			err = r.storage.Put(key, value, id, time.Now())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}(key, value)
	}

	c.JSON(http.StatusCreated, "OK")
	return
}
