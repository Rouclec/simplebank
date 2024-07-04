package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/rouclec/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(pool)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run a concurrent transfer transaction
	n := 5

	amount := 10.00

	errs := make(chan error)

	responses := make(chan TransfersTxResponse)

	for i := 0; i < n; i++ {
		go func() {
			response, err := store.TransferTx(context.Background(), TransferTxRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      "USD",
			})

			errs <- err
			responses <- response
		}()
	}

	// check responses
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		var fromAmount float64
		var toAmount float64
		var conversionError error

		err := <-errs

		require.NoError(t, err)

		response := <-responses

		require.NotEmpty(t, response)

		//check transfer
		transfer := response.Transfer

		fromAmount, conversionError = util.Converter(transfer.Currency, response.FromAccount.Currency, transfer.Amount)

		require.NoError(t, conversionError)

		toAmount, conversionError = util.Converter(transfer.Currency, response.ToAccount.Currency, transfer.Amount)

		require.NoError(t, conversionError)

		require.NotEmpty(t, transfer)

		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)

		require.NoError(t, err)

		// check entries
		fromEntry := response.FromEntry

		require.NotEmpty(t, fromEntry)

		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.NoError(t, err)
		require.Equal(t, fromEntry.Amount, -fromAmount)

		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)

		require.NoError(t, err)

		toEntry := response.ToEntry

		require.NotEmpty(t, toEntry)

		require.NoError(t, err)

		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, toAmount)

		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)

		require.NoError(t, err)

		//check accounts
		fromAcount := response.FromAccount

		require.NotEmpty(t, fromAcount)
		require.Equal(t, fromAcount.ID, account1.ID)

		toAcocunt := response.ToAccount

		require.NotEmpty(t, toAcocunt)
		require.Equal(t, toAcocunt.ID, account2.ID)

		diff := account1.Balance - fromAcount.Balance

		require.True(t, diff > 0)

		convertedDiff, err := util.Converter(response.FromAccount.Currency, "USD", diff)
		require.NoError(t, err)

		k := int(convertedDiff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)

		existed[k] = true
	}

	//check final updated balance of both accounts

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// require.Equal(t, updateAccount1.Balance, account1.Balance-float64(n)*amount)

	amountSent, err := util.Converter("USD", updateAccount1.Currency, (float64(n) * amount))

	require.NoError(t, err)

	require.Equal(t, fmt.Sprintf("%.2f", updateAccount1.Balance), fmt.Sprintf("%.2f", (account1.Balance-amountSent)))

	amountReceived, err := util.Converter("USD", updateAccount2.Currency, (float64(n) * amount))
	require.NoError(t, err)

	require.Equal(t, fmt.Sprintf("%.2f", updateAccount2.Balance), fmt.Sprintf("%.2f", (account2.Balance+amountReceived)))
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(pool)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	accountMap := make(map[int64]Accounts)

	// Add the accounts to the map
	accountMap[account1.ID] = account1
	accountMap[account2.ID] = account2

	// run a concurrent transfer transaction
	n := 10

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		amount := float64(10.00)

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
				Currency:      "USD",
			})

			errs <- err
		}()
	}

	// check responses
	for i := 0; i < n; i++ {
		err := <-errs

		require.NoError(t, err)
	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, updateAccount1.Balance, account1.Balance)

	require.Equal(t, updateAccount2.Balance, account2.Balance)
}
