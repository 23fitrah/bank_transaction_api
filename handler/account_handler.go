package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"transaction_api/constants"
	"transaction_api/dto/account"
	"transaction_api/model"
	"transaction_api/service"
	"transaction_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AccountHandler struct {
	accountService *service.AccountService
	validate       *validator.Validate
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		validate:       validator.New(),
	}
}

func (h *AccountHandler) GetInquiryAccountHandler(c *gin.Context) {
	var req account.Request

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

	data, err := h.accountService.GetInquiryAccountService(c, req.Request.ACCOUNT_NO)
	if err != nil {

		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedCheckAccountData,
		})
		return
	}

	if data == nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageNotFoundAccount,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, model.Response{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageSuccessGetAccount,
		Data:         data,
	})
}

func (h *AccountHandler) InsertAccountHandler(c *gin.Context) {
	var req account.AccountRequest

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

	reqData := req.Request

	count, err := h.accountService.CheckAccountNoExistsService(c, reqData.ACCOUNT_NO)
	if err != nil {

		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedCheckAccountData,
		})
		return
	}
	if count > 0 {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      "Account No " + reqData.ACCOUNT_NO + " already exists",
		})
		return
	}

	_, err = h.accountService.InsertAccountService(c, reqData)
	if err != nil {

		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedInsertAccount,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, model.Response{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageSuccessInsertAccount,
	})
}

func (h *AccountHandler) DeleteAccountHandler(c *gin.Context) {
	var req account.AccountDeleteRequest

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
	reqdata := req.Request
	count, err := h.accountService.CheckAccountIDExistsService(c, reqdata.ID_ACCOUNT)
	if err != nil {
		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedCheckAccountData,
		})
		return
	}
	if count == 0 {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      "Account ID " + fmt.Sprint(req.Request.ID_ACCOUNT) + " does not exist",
		})
		return
	}

	oldValue, err, isFound := h.getDetailAccount(c, reqdata.ID_ACCOUNT)
	if err != nil {
		log.Println("Failed to get detail account data:", err)
		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedGetAccount,
		})
		return
	}

	if !isFound {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageNotFoundAccount,
		})
		return
	}

	c.Set("oldValue", oldValue)

	deleteID, err := h.accountService.DeleteAccountService(c, req.Request.ID_ACCOUNT)
	if err != nil {

		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedDeleteAccount,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, model.Response{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageSuccessDeleteAccount,
		Data:         deleteID,
	})
}

func (h *AccountHandler) UpdateAccountHandler(c *gin.Context) {
	var req account.AccountRequest

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

	reqData := req.Request

	count, err := h.accountService.CheckAccountIDExistsService(c, reqData.ID_ACCOUNT)
	if err != nil {
		log.Println("Failed to check account data:", err)
		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedCheckAccountData,
		})
		return
	}
	if count == 0 {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      "Account ID " + fmt.Sprint(reqData.ID_ACCOUNT) + " does not exist",
		})
		return
	}

	oldValue, err, isFound := h.getDetailAccount(c, reqData.ID_ACCOUNT)
	if err != nil {
		log.Println("Failed to get detail account data:", err)
		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedGetAccount,
		})
		return
	}

	if !isFound {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageNotFoundAccount,
		})
		return
	}

	c.Set("oldValue", oldValue)

	_, err = h.accountService.UpdateAccountService(c, reqData)
	if err != nil {
		log.Println("Error updating account data:", err.Error())
		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedUpdateAccount,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, model.Response{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageSuccessUpdateAccount,
	})
}

func (h *AccountHandler) GetDetailAccountHandler(c *gin.Context) {
	idParam := c.Param("id")

	var id int64
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedBindJSON,
		})
		return
	}

	detail, err := h.accountService.GetDetailAccountService(c, id)
	if err != nil {
		log.Println("Error fetching account detail:", err)
		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedGetAccount,
		})
		return
	}

	if detail == nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageNotFoundAccount,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, model.Response{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageSuccessGetAccount,
		Data:         detail,
	})
}

func (h *AccountHandler) GetAllTransactionHandler(c *gin.Context) {
	var req account.AccountSearchRequest
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

	data, total, err := h.accountService.GetAllAccountsService(ctx, limit, offset, req.Request.ACCOUNT_NO, req.Request.NAME)
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
			Message:      constants.MessageNotFoundAccount,
		})
		return
	}

	pagination := utils.BuildPagination(page, limit, total)

	utils.WriteJSON(c, http.StatusOK, utils.PaginatedResponse{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageSuccessGetAccount,
		Data:         data,
		Pagination:   pagination,
	})
}

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
}
