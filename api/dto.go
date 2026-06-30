package api

import (
	"time"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
)

// AccountResponse is the client-facing representation of an account.
// Decoupling this from the db.Account model means adding columns to
// the database never leaks to API consumers.
type AccountResponse struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

// newAccountResponse maps a database Account into the API response shape.
func newAccountResponse(account db.Account) AccountResponse {
	return AccountResponse{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

// SuccessResponse provides a consistent shape for endpoints that don't
// return a resource body (e.g., DELETE).
type SuccessResponse struct {
	Message string `json:"message"`
}
