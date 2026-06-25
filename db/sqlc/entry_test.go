package db

import (
	"context"
	"testing"

	"github.com/NatdanaiKhe/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account *Account) *Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomInt(0, 1000),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	return &entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)
	require.NotEmpty(t, entry)
	require.Equal(t, account.ID, entry.AccountID)
	require.True(t, entry.Amount > 0)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	createdEntry := createRandomEntry(t, account)

	entry, err := testQueries.GetEntry(context.Background(), createdEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, createdEntry.ID, entry.ID)
	require.Equal(t, createdEntry.AccountID, entry.AccountID)
	require.Equal(t, createdEntry.Amount, entry.Amount)

}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	entries, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
		AccountID: account.ID,
		Limit:     10,
		Offset:    0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 10)
}
