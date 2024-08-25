package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kvgtl/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	senderAccount := createRandomAccount(t)
	receiverAccount := createRandomAccount(t)

	args := CreateTransferParams{
		FromAccountID: senderAccount.ID,
		ToAccountID:   receiverAccount.ID,
		Amount:        utils.RandomMoneyAmount(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, args.FromAccountID)
	require.Equal(t, transfer.ToAccountID, args.ToAccountID)
	require.Equal(t, transfer.Amount, args.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestDeleteTransfer(t *testing.T) {
	createdTransfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), createdTransfer.ID)

	require.NoError(t, err)

	deletedTransfer, err := testQueries.GetTransfer(context.Background(), createdTransfer.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedTransfer)
}

func TestGetTransfer(t *testing.T) {
	createdTransfer := createRandomTransfer(t)

	transfer, err := testQueries.GetTransfer(context.Background(), createdTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.ID, createdTransfer.ID)
	require.Equal(t, transfer.FromAccountID, createdTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, createdTransfer.ToAccountID)
	require.WithinDuration(t, transfer.CreatedAt, createdTransfer.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	args := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestUpdateTransfer(t *testing.T) {
	createdTransfer := createRandomTransfer(t)

	args := UpdateTransferParams{
		ID:          createdTransfer.ID,
		ToAccountID: createdTransfer.ToAccountID,
		Amount:      utils.RandomMoneyAmount(),
	}

	transfer, err := testQueries.UpdateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.ID, args.ID)
	require.Equal(t, transfer.FromAccountID, createdTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, args.ToAccountID)
	require.Equal(t, transfer.Amount, args.Amount)
	require.WithinDuration(t, transfer.CreatedAt, createdTransfer.CreatedAt, time.Second)
}
