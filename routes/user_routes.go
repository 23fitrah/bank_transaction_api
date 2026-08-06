package routes

import (
	"transaction_api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup, userHandler *handler.UserHandler) {
	r.POST("/login", userHandler.LoginUserHandler)
}
