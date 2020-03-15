package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Ping string `json:"ping"`
}

//-----------------------------------TEST GET
func PingGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}

//-----------------------------------TEST POST
func PingPost() gin.HandlerFunc {

	return func(c *gin.Context) {
		reqBody := Body{}
		rawData, _ := c.GetRawData()
		json.Unmarshal(rawData, &reqBody)
		c.JSON(http.StatusOK, gin.H{
			"posted": reqBody.Ping,
		})
	}
}
