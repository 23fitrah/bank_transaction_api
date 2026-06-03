package account

import (
	"transaction_api/dto"
)

type Request struct {
	dto.User
	Request RequestField `json:"request" validate:"required"`
}

type RequestField struct {
	ACCOUNT_NO string `json:"account_no" validate:"required,len=15"`
}

type AccountRequest struct {
	dto.User
	Request AccountRequestField `json:"request" validate:"required"`
}

type AccountRequestField struct {
	ID_ACCOUNT int64   `json:"id_account,omitempty"`
	NAME       string  `json:"account_name" validate:"required,max=50"`
	ACCOUNT_NO string  `json:"account_no" validate:"required,len=15"`
	ADDRESS    string  `json:"address" validate:"required,max=100"`
	COUNTRY    string  `json:"country" validate:"required,len=2"`
	EMAIL      string  `json:"email" validate:"required,email,max=50"`
	BALANCE    float64 `json:"balance" validate:"required,gt=0"`
	CURRENCY   string  `json:"currency" validate:"required,len=3"`
	STATUS     string  `json:"status" validate:"required,oneof=ACTIVE INACTIVE DORMANT FREEZE"`
}

type AccountDeleteRequest struct {
	dto.User
	Request AccountDeleteRequestField `json:"request" validate:"required"`
}

type AccountDeleteRequestField struct {
	ID_ACCOUNT int64 `json:"id_account" validate:"required"`
}

type AccountSearchRequest struct {
	dto.User
	Request AccountSearchRequestField `json:"request" validate:"required"`
}

type AccountSearchRequestField struct {
	ACCOUNT_NO string `json:"account_no,omitempty"`
	NAME       string `json:"account_name,omitempty"`
}
