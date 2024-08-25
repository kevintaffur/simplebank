package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kvgtl/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	return CreateRandomEntryAtSpecificAccount(t, account)
}

func CreateRandomEntryAtSpecificAccount(t *testing.T, account Account) Entry {
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoneyAmountForEntries(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.Amount, args.Amount)
	require.Equal(t, entry.AccountID, args.AccountID)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)

	require.NoError(t, err)

	entryToDelete, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entryToDelete)
}

func TestGetEntry(t *testing.T) {
	createdEntry := createRandomEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), createdEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.ID, createdEntry.ID)
	require.Equal(t, entry.AccountID, createdEntry.AccountID)
	require.Equal(t, entry.Amount, createdEntry.Amount)
	require.WithinDuration(t, entry.CreatedAt, createdEntry.CreatedAt, time.Second)

}

func TestUpdateEntry(t *testing.T) {
	createdEntry := createRandomEntry(t)

	args := UpdateEntryParams{
		ID:     createdEntry.ID,
		Amount: utils.RandomMoneyAmountForEntries(),
	}

	entry, err := testQueries.UpdateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.ID, args.ID)
	require.Equal(t, entry.AccountID, createdEntry.AccountID)
	require.Equal(t, entry.Amount, args.Amount)
	require.WithinDuration(t, entry.CreatedAt, createdEntry.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	args := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestAccountEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		CreateRandomEntryAtSpecificAccount(t, account)
	}

	args := ListAccountEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListAccountEntries(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
