package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	pool *pgxpool.Pool
}

type TransferTxRequest struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransfersTxResponse struct {
	Transfer    Transfers `json:"transfer"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

// CurrencyRate represents the exchange rate for a currency relative to USD
type CurrencyRate struct {
	Rate float64 `json:"rate"`
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		pool:    pool,
		Queries: New(pool),
	}
}

// Executes a function withing a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.pool.Begin(ctx)

	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

func converter(fromCurrency string, toCurrency string, amount float64) (float64, error) {

	rates := map[string]CurrencyRate{
		"EUR": {Rate: 1.1},    // Euros per USD
		"XAF": {Rate: 607.29}, // West African Francs per USD
		"CAD": {Rate: 1.35},   // Canadian dollar per USD
		"USD": {Rate: 1},      // USD per USD
	}

	if _, ok := rates[fromCurrency]; !ok {
		return 0, fmt.Errorf("unsupported currency: %s", fromCurrency)
	}
	if _, ok := rates[toCurrency]; !ok {
		return 0, fmt.Errorf("unsupported currency: %s", toCurrency)
	}

	if fromCurrency == toCurrency {
		return amount, nil
	}

	// Convert to USD first
	usdAmount := float64(amount) / rates[fromCurrency].Rate

	// Convert from USD to target currency
	convertedAmount := usdAmount * rates[toCurrency].Rate
	roundedAmount := fmt.Sprintf("%.2f", convertedAmount)
	parsedAmount, _ := strconv.ParseFloat(roundedAmount, 64)

	return parsedAmount, nil
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxRequest) (TransfersTxResponse, error) {
	var response TransfersTxResponse


	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var fromAccount Accounts
		var toAccount Accounts
		var convertedAmount float64

		if arg.FromAccountID < arg.ToAccountID {
			// Acquire locks on the accounts based on their IDs
			fromAccount, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil {
				return err
			}

			toAccount, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil {
				return err
			}

		} else {
			// Acquire locks on the accounts based on their IDs
			toAccount, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil {
				return err
			}

			fromAccount, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil {
				return err
			}

		}

		// Perform the transfer logic
		if balance := fromAccount.Balance; balance < arg.Amount {
			err = fmt.Errorf("insufficient balance to perform transaction. Balance is %v %v but attempted to transfer %v %v", balance, fromAccount.Currency, arg.Amount, fromAccount.Currency)
			return err
		}

		response.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
			FromCurrency:  fromAccount.Currency,
			ToCurrency:    toAccount.Currency,
		})
		if err != nil {
			return err
		}

		convertedAmount, err = converter(fromAccount.Currency, toAccount.Currency, arg.Amount)
		if err != nil {
			return err
		}

		response.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    convertedAmount,
		})
		if err != nil {
			return err
		}

		response.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update the balances of the accounts
		if fromAccount.ID < toAccount.ID {
			response.FromAccount, response.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, response.FromEntry.Amount, arg.ToAccountID, response.ToEntry.Amount)
		} else {
			response.ToAccount, response.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, response.ToEntry.Amount, arg.FromAccountID, response.FromEntry.Amount)
		}
		if err != nil {
			return err
		}

		return nil
	})

	return response, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 float64,
	accountID2 int64,
	amount2 float64,
) (account1 Accounts, account2 Accounts, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
