package transaction

type Transaction struct {
	ROWID_TRX         int64   `json:"rowid_trx"`
	ROWID_SENDER      int     `json:"rowid_sender"`
	ROWID_BENEFICIARY int     `json:"rowid_beneficiary"`
	TRANSACTION_DATE  string  `json:"transaction_date"`
	DEBET_ACCOUNT     string  `json:"debet_account"`
	DEBET_NAME        string  `json:"debet_name"`
	DEBET_CURR        string  `json:"debet_currency"`
	CREDIT_ACCOUNT    string  `json:"credit_account"`
	CREDIT_NAME       string  `json:"credit_name"`
	CREDIT_CURR       string  `json:"credit_currency"`
	AMOUNT            float64 `json:"amount"`
	REMARK            string  `json:"remark"`
	REFERENCE_NUMBER  string  `json:"reference_number"`
	STATUS_CODE       string  `json:"status_code"`
	STATUS_DESC       string  `json:"status_desc"`
}
