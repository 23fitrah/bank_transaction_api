package log_monitoring

import (
	"transaction_api/dto"
)

type RequestSearch struct {
	dto.User
	Request RequestSearchField `json:"request" validate:"required"`
}

type RequestSearchField struct {
	DATE_FROM string `json:"date_from" validate:"omitempty,datetime=2006-01-02"`
	DATE_TO   string `json:"date_to" validate:"omitempty,datetime=2006-01-02" `
	MENU      string `json:"menu,omitempty" `
	METHOD    string `json:"method" validate:"required"`
}
