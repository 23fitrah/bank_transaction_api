package routes

import (
	"transaction_api/handler"
	"transaction_api/utils"

	"github.com/gin-gonic/gin"
)

func RegisterTransactionRoutes(r *gin.Engine, transactionHandler *handler.TransactionHandler) {
	mx := r.Group("/api/v1/transaction")
	mx.Use(utils.AuthMiddlewareGin())
	{
		mx.POST("/create", transactionHandler.InsertTransactionHandler)
		mx.GET("/get-all", transactionHandler.GetAllTransactionHandler)
		mx.GET("/get-detail/:rowid_trx", transactionHandler.GetDetailTransactionHandler)
		mx.GET("/get-download", transactionHandler.GetDownloadTransactionHandler)
	}
}
