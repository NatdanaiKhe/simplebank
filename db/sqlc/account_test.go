package db

import (
	"context"
	"testing"

	"github.com/NatdanaiKhe/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) *Account {
	return createRandomAccountWithBalance(t, 0)
}

func createRandomAccountWithBalance(t *testing.T, balance int64) *Account {
	owner := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    owner,
		Balance:  balance,
		Currency: "THB",
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return &account
}

func createRandomUser(t *testing.T) string {
	username := util.RandomString(8)
	_, err := testDB.ExecContext(
		context.Background(),
		`INSERT INTO users (username, hashed_password, full_name, email)
		VALUES ($1, $2, $3, $4)`,
		username,
		util.RandomString(12),
		util.RandomString(12),
		username+"@example.com",
	)
	require.NoError(t, err)
	return username
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	result, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, account.ID, result.ID)
	require.Equal(t, account.Owner, result.Owner)
	require.Equal(t, account.Balance, result.Balance)
	require.Equal(t, account.Currency, result.Currency)
	require.Equal(t, account.CreatedAt, result.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: 200,
	}
	result := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, result)

	updatedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{Limit: 10, Offset: 0})
	require.Len(t, accounts, 10)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
	require.NoError(t, err)
	require.NotNil(t, accounts)
}

func TestListAccountsWithoutLimit(t *testing.T) {
	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{Limit: 0, Offset: 0})
	require.NoError(t, err)
	require.Nil(t, accounts)
}

func TestListAccountsWithNoArg(t *testing.T) {
	_, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{})
	require.Nil(t, err)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	_, err = testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
}
