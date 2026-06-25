package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides a transactional store for database operations.

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a transactional function with the store's database connection.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := store.Queries.WithTx(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Transfer transfers an amount from one account to another.
// It uses a transaction to ensure atomicity of the transfer operation.
func (store *Store) TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		if params.Amount <= 0 {
			return fmt.Errorf("amount must be positive")
		}

		if params.FromAccountID == params.ToAccountID {
			return fmt.Errorf("from_account_id and to_account_id must be different")
		}

		// Lock both way accounts for update
		// This condition is just deterministic ordering rule, so we lock the accounts in the same order in both branches
		// eg.
		// Transfer A -> B: lock 10, then 20
		// Transfer B -> A: lock 10, then 20
		var fromAccount Account
		if params.FromAccountID < params.ToAccountID {
			fromAccount, err = q.GetAccountForUpdate(ctx, params.FromAccountID)
			if err != nil {
				return err
			}
			_, err = q.GetAccountForUpdate(ctx, params.ToAccountID)
			if err != nil {
				return err
			}
		} else {
			_, err = q.GetAccountForUpdate(ctx, params.ToAccountID)
			if err != nil {
				return err
			}
			fromAccount, err = q.GetAccountForUpdate(ctx, params.FromAccountID)
			if err != nil {
				return err
			}
		}

		if fromAccount.Balance < params.Amount {
			return fmt.Errorf("insufficient balance")
		}
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID:   params.ToAccountID,
			Amount:        params.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount:    -params.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.ToAccountID,
			Amount:    params.Amount,
		})
		if err != nil {
			return err
		}
		// Update account balance
		// Subtract amount from sender, add amount to recipient
		result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:     params.FromAccountID,
			Amount: -params.Amount,
		})
		if err != nil {
			return err
		}

		// Add amount to recipient
		result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:     params.ToAccountID,
			Amount: params.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
