package service

import "errors"

var (
	ErrAccountNotFound      = errors.New("Account not found")
	ErrInsufficientFunds    = errors.New("Insufficient funds for transfer")
	ErrUnsupportedCurrency  = errors.New("Currency not supported")
	ErrAccountAlreadyExists = errors.New("Account already exists")
	ErrInternal             = errors.New("An internal error occurred")
)
