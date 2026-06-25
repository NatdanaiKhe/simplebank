package db

import (
	"context"
	"testing"

	"github.com/NatdanaiKhe/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) *Transfer {
	result, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomInt(1, 1000),
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, account1.ID, result.FromAccountID)
	require.Equal(t, account2.ID, result.ToAccountID)
	require.NotZero(t, result.Amount)
	return &result
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transfer := createRandomTransfer(t, *account1, *account2)
	require.NotEmpty(t, transfer)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transfer := createRandomTransfer(t, *account1, *account2)
	require.NotEmpty(t, transfer)

	result, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, transfer.ID, result.ID)
	require.Equal(t, account1.ID, result.FromAccountID)
	require.Equal(t, account2.ID, result.ToAccountID)
	require.NotZero(t, result.Amount)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transfers := make([]*Transfer, 0, 10)

	for i := 0; i < 10; i++ {
		transfer := createRandomTransfer(t, *account1, *account2)
		transfers = append(transfers, transfer)
	}

	result, err := testQueries.ListTransfers(context.Background(), ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         10,
		Offset:        0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Len(t, result, 10)
	for i, transfer := range result {
		require.Equal(t, transfers[i].ID, transfer.ID)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.Amount)
	}
}
