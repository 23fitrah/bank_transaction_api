package account

type Account struct {
	ID_ACCOUNT int64   `json:"id_account,omitempty"`
	NAME       string  `json:"account_name"`
	ACCOUNT_NO string  `json:"account_no"`
	ADDRESS    string  `json:"address,omitempty"`
	COUNTRY    string  `json:"country,omitempty"`
	EMAIL      string  `json:"email,omitempty"`
	BALANCE    float64 `json:"balance"`
	CURRENCY   string  `json:"currency"`
	STATUS     string  `json:"status"`
}
