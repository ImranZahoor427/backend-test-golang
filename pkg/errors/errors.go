package errors

import "errors"

var (
	// Common errors
	ErrInvalidInput           = errors.New("invalid input")
	ErrAccountNotFound        = errors.New("account not found")
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrTransactionFailed      = errors.New("transaction failed")
	ErrDuplicateRequest       = errors.New("duplicate request")
	ErrInvalidAmount          = errors.New("amount must be positive")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
	ErrInvalidLimit           = errors.New("invalid limit")
	ErrInvalidOffset          = errors.New("invalid offset")
	ErrParsingID              = errors.New("failed to parse ID")
	ErrFetchingTransactions   = errors.New("failed to retrieve transactions")
	ErrFake                   = errors.New("fake error")
)
