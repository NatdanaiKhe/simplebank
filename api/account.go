package api

import (
	"net/http"

	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/gin-gonic/gin"
)

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

func (server *Server) createAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, err)
		return
	}

	account, err := server.service.Create(c, service.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  req.Balance,
		Currency: req.Currency,
	})
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, newAccountResponse(account))
}

func (server *Server) getAccount(c *gin.Context) {
	var req GetAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errorResponse(c, err)
		return
	}

	account, err := server.service.GetByID(c, req.ID)
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, newAccountResponse(account))
}

func (server *Server) listAccounts(c *gin.Context) {
	var param ListAccountsRequest
	if err := c.ShouldBindQuery(&param); err != nil {
		errorResponse(c, err)
		return
	}

	accounts, total, err := server.service.List(c, service.ListAccountsParams{
		Limit:  param.PageSize,
		Offset: param.PageSize * (param.PageNumber - 1),
	})
	if err != nil {
		errorResponse(c, err)
		return
	}

	accountResponses := make([]AccountResponse, len(accounts))
	for i, a := range accounts {
		accountResponses[i] = newAccountResponse(a)
	}

	res := ListAccountsResponse{
		Accounts:   accountResponses,
		PageNumber: param.PageNumber,
		PageSize:   param.PageSize,
		Total:      int32(total),
	}
	c.JSON(http.StatusOK, res)
}

func (server *Server) updateAccount(c *gin.Context) {
	var uri UpdateAccountUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorResponse(c, err)
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, err)
		return
	}

	account, err := server.service.Update(c, service.UpdateAccountParams{
		ID:      uri.ID,
		Balance: req.Balance,
	})
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, newAccountResponse(account))
}

func (server *Server) deleteAccount(c *gin.Context) {
	var uri DeleteAccountUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorResponse(c, err)
		return
	}

	err := server.service.Delete(c, uri.ID)
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Message: "account deleted"})
}
