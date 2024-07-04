package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rouclec/simplebank/util"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxRequest) (TransfersTxResponse, error)
}

type SQLStore struct {
	*Queries
	pool *pgxpool.Pool
}

type TransferTxRequest struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
}

type TransfersTxResponse struct {
	Transfer    Transfers `json:"transfer"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

// CurrencyRate represents the exchange rate for a currency relative to USD

func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		pool:    pool,
		Queries: New(pool),
	}
}

// Executes a function withing a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxRequest) (TransfersTxResponse, error) {
	var response TransfersTxResponse

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var fromAccount Accounts
		var toAccount Accounts
		var fromAmount float64
		var toAmount float64

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
		fromAmount, err = util.Converter(arg.Currency, fromAccount.Currency, arg.Amount)
		if err != nil {
			return err
		}

		toAmount, err = util.Converter(arg.Currency, toAccount.Currency, arg.Amount)
		if err != nil {
			return err
		}

		if balance := fromAccount.Balance; balance < fromAmount {
			err = fmt.Errorf("insufficient balance to perform transaction")
			return err
		}

		response.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
			Currency:      arg.Currency,
		})
		if err != nil {
			return err
		}

		response.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    toAmount,
		})
		if err != nil {
			return err
		}

		response.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -fromAmount,
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
