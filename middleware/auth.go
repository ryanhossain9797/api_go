package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CheckAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("authorized")
	}
}
