package routes

import (
	"transaction_api/handler"
	"transaction_api/utils"

	"github.com/gin-gonic/gin"
)

func RegisterAccountRoutes(r *gin.Engine, accountHandler *handler.AccountHandler) {
	mx := r.Group("/api/v1/account")
	mx.Use(utils.AuthMiddlewareGin())
	{
		mx.POST("/inquiry", accountHandler.GetInquiryAccountHandler)
		mx.GET("/get-all", accountHandler.GetAllTransactionHandler)
		mx.PUT("/update", accountHandler.UpdateAccountHandler)
		mx.DELETE("/delete", accountHandler.DeleteAccountHandler)
		mx.GET("/get-detail/:id", accountHandler.GetDetailAccountHandler)
		mx.POST("/create", accountHandler.InsertAccountHandler)
	}
}
