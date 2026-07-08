package service

import (
	"context"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
)

type TransferService interface {
	CreateTransfer(ctx context.Context, param db.CreateTransferParams) error
}
type transferService struct {
	store db.Store
}

func NewTransferService(store db.Store) TransferService {
	return transferService{store: store}
}

func (s transferService) CreateTransfer(ctx context.Context, param db.CreateTransferParams) error {
	transferTxParams := db.TransferTxParams{
		FromAccountID: param.FromAccountID,
		ToAccountID:   param.ToAccountID,
		Amount:        param.Amount,
	}
	_, err := s.store.TransferTx(ctx, transferTxParams)
	return err
}
