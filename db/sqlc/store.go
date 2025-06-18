package db

import (
	"context"
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

func (store *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx) // rollback the transaction if it hasn't been committed

	q := store.WithTx(tx)

	err = fn(q)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}