package db

import (
	"context"
	"testing"
	"time"

	"github.com/rouclec/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Accounts {
	currency := util.RandomCurrency()

	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(currency),
		Currency: currency,
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, account.Owner, arg.Owner)
	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.Currency, arg.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	accountFound, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)

	require.Equal(t, accountFound.Balance, account.Balance)
	require.Equal(t, accountFound.Currency, account.Currency)
	require.Equal(t, accountFound.Owner, account.Owner)

	require.WithinDuration(t, accountFound.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: account.Balance + util.RandomBalance(account.Currency),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, updatedAccount)

	require.Equal(t, updatedAccount.ID, account.ID)
	require.Equal(t, updatedAccount.Balance, arg.Balance)
	require.Equal(t, updatedAccount.Currency, account.Currency)
	require.Equal(t, updatedAccount.Owner, account.Owner)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

	accountFound, err1 := testQueries.GetAccount(context.Background(), account.ID)

	require.Error(t, err1)

	require.EqualError(t, err1, ErrRecordNotFound.Error())
	require.Empty(t, accountFound)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
