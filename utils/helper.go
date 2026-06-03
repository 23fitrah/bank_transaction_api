package utils

import (
	"github.com/gin-gonic/gin"
)

func WriteJSON(c *gin.Context, statusCode int, data interface{}) {
	c.Set("response_msg", data)
	c.JSON(statusCode, data)
}
