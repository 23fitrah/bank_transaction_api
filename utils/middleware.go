package utils

import (
	"log"
	"net/http"
	"os"
	"strings"
	"transaction_api/model"

	"github.com/gin-gonic/gin"
)

func AuthMiddlewareGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("REQ PATH:", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.JSON(http.StatusUnauthorized, model.Response{
				ResponseCode: "0002",
				Message:      "Unauthorized",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)
		authToken := os.Getenv("TOKEN_AUTH")

		if token != authToken {
			c.JSON(http.StatusUnauthorized, model.Response{
				ResponseCode: "0002",
				Message:      "Invalid token",
			})
			c.Abort()
			return
		}

		// Lanjut ke handler berikutnya
		c.Next()
	}
}

func Log() {

}
