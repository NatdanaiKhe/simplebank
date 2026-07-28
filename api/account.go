package api

import (
	"errors"
	"net/http"

	"github.com/NatdanaiKhe/simplebank/api/dto"
	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/gin-gonic/gin"
)

func (server *Server) createAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, err)
		return
	}

	params := service.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  req.Balance,
		Currency: req.Currency,
	}
	account, err := server.services.AccountService.Create(c, params)

	if err != nil {
		if errors.Is(err, service.ErrForeignKeyViolation) {
			errorResponse(c, service.ErrForeignKeyViolation)
			return
		}

		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewAccountResponse(account))
}

func (server *Server) getAccount(c *gin.Context) {
	var req dto.GetAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errorResponse(c, err)
		return
	}

	account, err := server.services.AccountService.GetByID(c, req.ID)
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewAccountResponse(account))
}

func (server *Server) listAccounts(c *gin.Context) {
	var param dto.ListAccountsRequest
	if err := c.ShouldBindQuery(&param); err != nil {
		errorResponse(c, err)
		return
	}

	accounts, total, err := server.services.AccountService.List(c, service.ListAccountsParams{
		Limit:  param.PageSize,
		Offset: param.PageSize * (param.PageNumber - 1),
	})
	if err != nil {
		errorResponse(c, err)
		return
	}

	accountResponses := make([]dto.AccountResponse, len(accounts))
	for i, a := range accounts {
		accountResponses[i] = dto.NewAccountResponse(a)
	}

	res := dto.ListAccountsResponse{
		Accounts:   accountResponses,
		PageNumber: param.PageNumber,
		PageSize:   param.PageSize,
		Total:      int32(total),
	}
	c.JSON(http.StatusOK, res)
}

func (server *Server) updateAccount(c *gin.Context) {
	var uri dto.UpdateAccountUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorResponse(c, err)
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, err)
		return
	}

	account, err := server.services.AccountService.Update(c, service.UpdateAccountParams{
		ID:      uri.ID,
		Balance: req.Balance,
	})
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewAccountResponse(account))
}

func (server *Server) deleteAccount(c *gin.Context) {
	var uri dto.DeleteAccountUri
	if err := c.ShouldBindUri(&uri); err != nil {
		errorResponse(c, err)
		return
	}

	err := server.services.AccountService.Delete(c, uri.ID)
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "account deleted"})
}
