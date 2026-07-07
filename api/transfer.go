package api

import (
	"net/http"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR THB"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, err)
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result := server.service.TransferService.CreateTransfer(ctx, arg)
	if result != nil {
		errorResponse(ctx, result)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func validateTransferRequest(ctx *gin.Context, server *Server, req *TransferRequest, accountID int) bool {
	account, err := server.service.AccountService.GetByID(ctx, int64(accountID))
	if err != nil {
		return false
	}

	if req.FromAccountID <= 0 || req.ToAccountID <= 0 || req.Amount <= 0 {
		return false
	}
	if req.FromAccountID != int64(accountID) && req.ToAccountID != int64(accountID) {
		return false
	}
	if req.Currency != account.Currency {
		return false
	}

	return true
}
