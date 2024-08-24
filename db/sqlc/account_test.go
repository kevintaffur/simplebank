package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kevtl/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoneyAmount(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), createdAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.ID, createdAccount.ID)
	require.Equal(t, account.Owner, createdAccount.Owner)
	require.Equal(t, account.Balance, createdAccount.Balance)
	require.Equal(t, account.Currency, createdAccount.Currency)
	require.WithinDuration(t, account.CreatedAt, createdAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	newBalance := utils.RandomMoneyAmount()

	args := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: newBalance,
	}

	account, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.ID, args.ID)
	require.Equal(t, account.Owner, createdAccount.Owner)
	require.Equal(t, account.Currency, createdAccount.Currency)
	require.WithinDuration(t, account.CreatedAt, createdAccount.CreatedAt, time.Second)
	require.Equal(t, account.Balance, args.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

	accountDeleted, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountDeleted)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
