package transaction

import (
	"transaction_api/dto"
)

type RequestSearch struct {
	dto.User
	Request RequestSearchField `json:"request" validate:"required"`
}

type RequestSearchField struct {
	DATE_FROM        string `json:"date_from" validate:"omitempty,datetime=2006-01-02"`
	DATE_TO          string `json:"date_to" validate:"omitempty,datetime=2006-01-02" `
	REFERENCE_NUMBER string `json:"reference_number,omitempty" `
}

type Request struct {
	dto.User
	Request RequestField `json:"request" validate:"required"`
}

type RequestField struct {
	ROWID_TRX         int64   `json:"rowid_trx,omitempty" `
	ROWID_SENDER      int     `json:"rowid_sender" validate:"required"`
	ROWID_BENEFICIARY int     `json:"rowid_beneficiary" validate:"required"`
	TRANSACTION_DATE  string  `json:"transaction_date,omitempty"`
	REFERENCE_NUMBER  string  `json:"reference_number,omitempty"`
	DEBET_ACCOUNT     string  `json:"debet_account" validate:"required,len=15"`
	DEBET_NAME        string  `json:"debet_name" validate:"required,max=50"`
	DEBET_CURR        string  `json:"debet_curr" validate:"required,len=3"`
	CREDIT_ACCOUNT    string  `json:"credit_account" validate:"required,len=15"`
	CREDIT_NAME       string  `json:"credit_name" validate:"required,max=50"`
	CREDIT_CURR       string  `json:"credit_curr" validate:"required,len=3"`
	AMOUNT            float64 `json:"amount" validate:"required"`
	REMARK            string  `json:"remark" validate:"required,max=100"`
}
