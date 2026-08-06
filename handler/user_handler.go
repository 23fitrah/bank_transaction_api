package handler

import (
	"net/http"
	"transaction_api/constants"
	"transaction_api/dto/user"
	"transaction_api/model"
	"transaction_api/service"
	"transaction_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService *service.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

func (h *UserHandler) LoginUserHandler(c *gin.Context) {
	var req user.Request

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

	data, err := h.userService.GetUserService(c, req.Username, req.Password)
	if err != nil {

		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      err.Error(),
		})
		return
	}

	if data == nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageUserNotFound,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, model.Response{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageUserFound,
		Data:         data,
	})
}
