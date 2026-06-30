package api

import (
	"log"
	"net/http"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Balance  int64  `json:"balance" binding:"required"`
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
	Accounts   []db.Account `json:"accounts"`
	PageNumber int32        `json:"page_number"`
	PageSize   int32        `json:"page_size"`
	Total      int32        `json:"total"`
}

type UpdateAccountRequest struct {
	ID      int64 `json:"id" binding:"required,min=1"`
	Balance int64 `json:"balance" binding:"required"`
}

type DeleteAccountRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (server *Server) createAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, err)
		return
	}

	args := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  req.Balance,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(c, args)
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, account)
}

func (server *Server) getAccount(c *gin.Context) {
	var req GetAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errorResponse(c, err)
		return
	}

	account, err := server.store.GetAccount(c, req.ID)
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, account)
}

func (server *Server) listAccounts(c *gin.Context) {
	var param ListAccountsRequest
	log.Println("listAccounts", param)
	if err := c.ShouldBindQuery(&param); err != nil {
		errorResponse(c, err)
		return
	}

	args := db.ListAccountsParams{
		Limit:  param.PageSize,
		Offset: param.PageSize * (param.PageNumber - 1),
	}
	accounts, err := server.store.ListAccounts(c, args)
	if err != nil {
		errorResponse(c, err)
		return
	}

	accountsCount, err := server.store.CountAccounts(c)
	if err != nil {
		errorResponse(c, err)
		return
	}

	res := ListAccountsResponse{
		Accounts:   accounts,
		PageNumber: param.PageNumber,
		PageSize:   param.PageSize,
		Total:      int32(accountsCount),
	}
	c.JSON(http.StatusOK, res)
}

func (server *Server) updateAccount(c *gin.Context) {
	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, err)
		return
	}

	args := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}

	err := server.store.UpdateAccount(c, args)
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, args)
}

func (server *Server) deleteAccount(c *gin.Context) {
	var req DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, err)
		return
	}

	err := server.store.DeleteAccount(c, req.ID)
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}
