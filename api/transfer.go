package api

import (
	"net/http"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, err)
		return
	}

	if !server.validateTransferRequest(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validateTransferRequest(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		errorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateTransferRequest(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.services.AccountService.GetByID(ctx, accountID)
	if err != nil {
		return false
	}

	if currency != account.Currency {
		errorResponse(ctx, service.ErrTransferCurrencyMismatch)
		return false
	}

	return true
}
