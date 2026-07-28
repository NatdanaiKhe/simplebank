package dto

import (
	"time"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
)

// Input DTOs
type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required,min=3,max=50"`
	Balance  int64  `json:"balance" binding:"required,min=0"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR THB"`
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListAccountsRequest struct {
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=1,max=10"`
}

type ListAccountsResponse struct {
	Accounts   []AccountResponse `json:"accounts"`
	PageNumber int32             `json:"page_number"`
	PageSize   int32             `json:"page_size"`
	Total      int32             `json:"total"`
}

type UpdateAccountRequest struct {
	Balance int64 `json:"balance" binding:"required,min=0"`
}

type UpdateAccountUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type DeleteAccountUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// Output DTOs
type AccountResponse struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccountResponse(account db.Account) AccountResponse {
	return AccountResponse{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func NewSuccessResponse(message string) SuccessResponse {
	return SuccessResponse{Message: message}
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserResponse(user db.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
