package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute SQL queries and transactions
type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) ExecTx(ctx context.Context, fn func(*Queries) error) (err error) {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed to rollback transaction: %v, original error: %v", rbErr, err)
			}
		}
	}()

	q := store.WithTx(tx)

	err = fn(q)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer     Transfer `json:"transfer"`
	FromAccount  Account  `json:"from_account"`
	ToAccount    Account  `json:"to_account"`
	FromEntry    Entry    `json:"from_entry"`
	ToEntry      Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	
	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error
		// Create transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		
		// Create from entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create to entry
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		
		// Update account balances
		result.FromAccount, result.ToAccount, err = updateBalances(
			ctx, q, arg.FromAccountID, arg.ToAccountID, arg.Amount,
		)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func updateBalances(
	ctx context.Context,
	q *Queries,
	fromID, toID int64,
	amount int64,
) (fromAccount, toAccount Account, err error) {
	// сортировка по ID для предотвращения дедлоков
	if fromID < toID {
		fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     fromID,
			Amount: -amount,
		})
		if err != nil {
			return
		}
		toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     toID,
			Amount: amount,
		})
	} else {
		toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     toID,
			Amount: amount,
		})
		if err != nil {
			return
		}
		fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     fromID,
			Amount: -amount,
		})
	}
	return
}
