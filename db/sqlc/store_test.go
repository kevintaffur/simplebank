package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	senderAccount := createRandomAccount(t)
	receiverAccount := createRandomAccount(t)

	fmt.Println("before transaction>>Sender:", senderAccount.Balance)
	fmt.Println("before transaction>>Receiver:", receiverAccount.Balance)

	// run n concurrent transfer transactions.
	n := 5
	amount := int64(10)

	// channels are used for  communication between goroutines
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: senderAccount.ID,
				ToAccountID:   receiverAccount.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check the results of goroutines.
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)

		require.Equal(t, senderAccount.ID, transfer.FromAccountID)
		require.Equal(t, receiverAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries.
		// Sender's entry.
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, senderAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// Receiver's entry.
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, receiverAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, senderAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, receiverAccount.ID, toAccount.ID)

		fmt.Println("transaction>>Sender:", fromAccount.Balance)
		fmt.Println("transaction>>Receiver:", toAccount.Balance)

		// check accounts' balance.
		senderDifference := senderAccount.Balance - fromAccount.Balance
		receiverDifference := toAccount.Balance - receiverAccount.Balance

		fmt.Println("difference>>Sender:", senderDifference)
		fmt.Println("difference>>Receiver:", receiverDifference)

		require.Equal(t, senderDifference, receiverDifference)
		require.True(t, senderDifference > 0) // checking one of them is enough, both should be the same.
		require.True(t, senderDifference%amount == 0)

		k := int(senderDifference / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final updated balances.
	updatedSenderAccount, err := testQueries.GetAccount(context.Background(), senderAccount.ID)
	require.NoError(t, err)

	updatedReceiverAccount, err := testQueries.GetAccount(context.Background(), receiverAccount.ID)
	require.NoError(t, err)

	fmt.Println("after transaction>>Sender:", updatedSenderAccount.Balance)
	fmt.Println("after transaction>>Receiver:", updatedReceiverAccount.Balance)

	require.Equal(t, senderAccount.Balance-int64(n)*amount, updatedSenderAccount.Balance)
	require.Equal(t, receiverAccount.Balance+int64(n)*amount, updatedReceiverAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	senderAccount := createRandomAccount(t)
	receiverAccount := createRandomAccount(t)

	// run n concurrent transfer transactions.
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := senderAccount.ID
		toAccountID := receiverAccount.ID

		if i%2 == 1 {
			fromAccountID = receiverAccount.ID
			toAccountID = senderAccount.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check final updated balances.
	updatedSenderAccount, err := testQueries.GetAccount(context.Background(), senderAccount.ID)
	require.NoError(t, err)

	updatedReceiverAccount, err := testQueries.GetAccount(context.Background(), receiverAccount.ID)
	require.NoError(t, err)

	require.Equal(t, senderAccount.Balance, updatedSenderAccount.Balance)
	require.Equal(t, receiverAccount.Balance, updatedReceiverAccount.Balance)
}
