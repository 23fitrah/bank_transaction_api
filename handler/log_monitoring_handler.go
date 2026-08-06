package handler

import (
	"log"
	"net/http"
	"transaction_api/constants"
	"transaction_api/dto/log_monitoring"
	"transaction_api/model"
	"transaction_api/service"
	"transaction_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LogHandler struct {
	logService *service.LogMonitoringService
	validate   *validator.Validate
}

func NewLogHandler(logService *service.LogMonitoringService) *LogHandler {
	return &LogHandler{
		logService: logService,
		validate:   validator.New(),
	}
}

func (h *LogHandler) GetAllLogMonitoringHandler(c *gin.Context) {
	var req log_monitoring.RequestSearch
	ctx := c.Request.Context()
	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedBindJSON,
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      utils.FormatValidationErrors(err)["reason"],
		})
		return
	}

	data, total, err := h.logService.GetAllLogService(ctx, limit, offset, req.Request.MENU, req.Request.DATE_FROM, req.Request.DATE_TO, req.Request.METHOD)
	if err != nil {
		log.Println("Failed to fetch account data", err)
		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedGetAccount,
		})
		return
	}
	if total <= 0 {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageNotFoundLog,
		})
		return
	}

	pagination := utils.BuildPagination(page, limit, total)

	utils.WriteJSON(c, http.StatusOK, utils.PaginatedResponse{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageSuccessGetLog,
		Data:         data,
		Pagination:   pagination,
	})
}

/*
func (h *AccountHandler) getDetailAccount(c *gin.Context, id int64) (string, error, bool) {
	detail, err := h.accountService.GetDetailAccountService(c.Request.Context(), id)

	if err != nil {
		return "", err, false
	}

	if detail == nil {
		return "", nil, false
	}

	detailJSON, _ := json.Marshal(detail)
	return string(detailJSON), nil, true
}*/
