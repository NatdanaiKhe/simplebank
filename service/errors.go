package service

import "errors"

var (
	// Account Errors
	ErrAccountNotFound      = errors.New("Account not found")
	ErrInsufficientFunds    = errors.New("Insufficient funds for transfer")
	ErrUnsupportedCurrency  = errors.New("Currency not supported")
	ErrAccountAlreadyExists = errors.New("Account already exists")
	ErrInternal             = errors.New("An internal error occurred")

	// Transfer Errors
	ErrTransferSameAccount       = errors.New("Transfer same account")
	ErrTransferCurrencyMismatch  = errors.New("Transfer currency mismatch")
	ErrTransferAmountInvalid     = errors.New("Transfer amount invalid")
	ErrTransferInsufficientFunds = errors.New("Transfer insufficient funds")
)
