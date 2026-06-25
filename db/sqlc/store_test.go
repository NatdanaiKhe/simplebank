package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	n := 5
	amount := int64(100)
	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	senderAccount := createRandomAccountWithBalance(t, amount*int64(n))
	recipientAccount := createRandomAccount(t)
	// run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: senderAccount.ID,
				ToAccountID:   recipientAccount.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool, n)
	// wait for all transactions to complete
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		// verify the transfer was created and can be retrieved
		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.Equal(t, senderAccount.ID, transfer.FromAccountID)
		require.Equal(t, recipientAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from entry was created
		fromEntry := result.FromEntry
		require.Equal(t, senderAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check to entry was created
		toEntry := result.ToEntry
		require.Equal(t, recipientAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, senderAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, recipientAccount.ID, toAccount.ID)

		// check balances were updated
		diffBalance := senderAccount.Balance - fromAccount.Balance
		diffBalanceUpdated := toAccount.Balance - recipientAccount.Balance

		require.Equal(t, diffBalance, diffBalanceUpdated)
		require.True(t, diffBalance > 0)
		require.True(t, diffBalanceUpdated > 0)
		require.True(t, diffBalance%amount == 0) // 1 * amount, 2 * amount, ...

		k := int(diffBalance / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check all transfers were processed
	updatedSenderAccount, err := store.GetAccount(context.Background(), senderAccount.ID)
	require.NoError(t, err)
	// updated sender account balance should be align with transfer amount, i.e. initial balance - n * amount
	require.Equal(t, senderAccount.Balance-int64(n)*amount, updatedSenderAccount.Balance)

	updatedRecipientAccount, err := store.GetAccount(context.Background(), recipientAccount.ID)
	require.NoError(t, err)
	// updated recipient account balance should be align with transfer amount, i.e. initial balance + n * amount
	require.Equal(t, recipientAccount.Balance+int64(n)*amount, updatedRecipientAccount.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	n := 10
	amount := int64(10)
	initialBalance := amount * int64(n)

	account1 := createRandomAccountWithBalance(t, initialBalance)
	account2 := createRandomAccountWithBalance(t, initialBalance)

	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		require.Equal(t, amount, result.Transfer.Amount)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.Equal(t, amount, result.ToEntry.Amount)
	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, initialBalance, updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, initialBalance, updatedAccount2.Balance)
}
