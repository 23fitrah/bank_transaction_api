package constants

const (
	// Cover / Tracer
	CodeSuccess = "00"
	CodeFailed  = "01"

	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"

	MessageFailedBindJSON = "Failed to bind JSON"

	//account
	MessageSuccessInsertAccount   = "Account successfully inserted"
	MessageFailedInsertAccount    = "Failed to insert account"
	MessageSuccessUpdateAccount   = "Account successfully updated"
	MessageFailedUpdateAccount    = "Failed to update account"
	MessageSuccessGetAccount      = "Account successfully retrieved"
	MessageFailedGetAccount       = "Failed to retrieve account"
	MessageSuccessDeleteAccount   = "Account successfully deleted"
	MessageFailedDeleteAccount    = "Failed to delete account"
	MessageFailedCheckAccountData = "Failed to check account data"
	MessageNotFoundAccount        = "Account not found"

	//transaction
	MessageSuccessInsertTransaction  = "Transaction successfully inserted"
	MessageFailedInsertTransaction   = "Failed to insert transaction"
	MessageSuccessGetTransaction     = "Transaction successfully retrieved"
	MessageFailedGetTransaction      = "Failed to retrieve transaction"
	MessageSuccessUpdateTransaction  = "Transaction successfully updated"
	MessageFailedUpdateTransaction   = "Failed to update transaction"
	MessageSuccessDeleteTransaction  = "Transaction successfully deleted"
	MessageFailedDeleteTransaction   = "Failed to delete transaction"
	MessageFailedCheckBalance        = "Failed to check account balance"
	MessageFailedInsufficientBalance = "Insufficient balance"
	MMessageDataNotFoundTransaction  = "Data transaction not found"
)
