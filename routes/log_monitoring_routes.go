package routes

import (
	"transaction_api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterLogMonitoringRoutes(r *gin.Engine, logMonitoringHandler *handler.LogHandler) {
	mx := r.Group("/api/v1/log")
	//mx.Use(utils.AuthMiddlewareGin())
	//{
	mx.POST("/get-all", logMonitoringHandler.GetAllLogMonitoringHandler)

	//}
}
