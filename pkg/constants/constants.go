package constants

// TransactionType is a deposit or withdrawal
type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
)
