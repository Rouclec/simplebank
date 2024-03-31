package db

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(pool)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run a concurrent transfer transaction
	n := 5
	var amount float64

	amounts := map[string]CurrencyRate{
		"EUR": {Rate: 1100},  // Euros per USD
		"XAF": {Rate: 60729}, // West African Francs per USD
		"CAD": {Rate: 135},   // Canadian dollar per USD
		"USD": {Rate: 100},   // USD per USD
	}


	amount = amounts[account1.Currency].Rate

	errs := make(chan error)

	responses := make(chan TransfersTxResponse)

	for i := 0; i < n; i++ {
		go func() {
			response, err := store.TransferTx(context.Background(), TransferTxRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			responses <- response
		}()
	}

	// check responses
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs

		require.NoError(t, err)

		response := <-responses

		require.NotEmpty(t, response)

		//check transfer
		transfer := response.Transfer

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
		require.Equal(t, fromEntry.Amount, -amount)

		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)

		require.NoError(t, err)

		toEntry := response.ToEntry

		require.NotEmpty(t, toEntry)

		toAmount, err := converter(transfer.FromCurrency, transfer.ToCurrency, amount)

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

		k := int(diff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)

		existed[k] = true
	}

	//check final updated balance of both accounts
	tolerance := 0.001 // Adjust tolerance as needed

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// require.Equal(t, updateAccount1.Balance, account1.Balance-float64(n)*amount)
	require.True(t, math.Abs(updateAccount1.Balance-(account1.Balance-float64(n)*amount)) <= tolerance)

	amountReceived, _ := converter(updateAccount1.Currency, updateAccount2.Currency, amount)
	// require.Equal(t, updateAccount2.Balance, account2.Balance+float64(n)*amountReceived)
	require.True(t, math.Abs(updateAccount2.Balance-(account2.Balance+float64(n)*amountReceived)) <= tolerance)
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

	amounts := map[string]CurrencyRate{
		"EUR": {Rate: 11.0},   // Euros per USD
		"XAF": {Rate: 6072.9}, // West African Francs per USD
		"CAD": {Rate: 13.5},   // Canadian dollar per USD
		"USD": {Rate: 10.0},   // USD per USD
	}

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		fromAccount := accountMap[fromAccountID]

		amount := amounts[fromAccount.Currency].Rate

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
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
