package service

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/lib/pq"
)

// AccountService defines the business operations for accounts.
// It sits between the HTTP layer and the data layer, keeping
// business logic out of handlers.
type AccountService interface {
	Create(ctx context.Context, params CreateAccountParams) (db.Account, error)
	GetByID(ctx context.Context, id int64) (db.Account, error)
	List(ctx context.Context, params ListAccountsParams) ([]db.Account, int64, error)
	Update(ctx context.Context, params UpdateAccountParams) (db.Account, error)
	Delete(ctx context.Context, id int64) error
}

// Service-specific parameter types. These look similar to the db package
// types today, but they can diverge as business rules evolve without
// changing the handler signatures or the store interface.

type CreateAccountParams struct {
	Owner    string
	Balance  int64
	Currency string
}

type ListAccountsParams struct {
	Limit  int32
	Offset int32
}

type UpdateAccountParams struct {
	ID      int64
	Balance int64
}

// accountService is the production implementation backed by db.Store.
type accountService struct {
	store db.Store
}

// NewAccountService wires the production service to a store.
func NewAccountService(store db.Store) AccountService {
	return &accountService{store: store}
}

func (s *accountService) Create(ctx context.Context, params CreateAccountParams) (db.Account, error) {
	dbParams := db.CreateAccountParams{
		Owner:    params.Owner,
		Balance:  params.Balance,
		Currency: params.Currency,
	}
	account, err := s.store.CreateAccount(ctx, dbParams)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return db.Account{}, ErrAccountAlreadyExists
		}
		return db.Account{}, ErrInternal
	}
	return account, nil
}

func (s *accountService) GetByID(ctx context.Context, id int64) (db.Account, error) {
	account, err := s.store.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Account{}, ErrAccountNotFound
		}
		return db.Account{}, ErrInternal
	}
	return account, nil
}

// List returns the page of accounts plus the total count in the database.
// Bundling these into a single call keeps pagination logic out of handlers.
func (s *accountService) List(ctx context.Context, params ListAccountsParams) ([]db.Account, int64, error) {
	accounts, err := s.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := s.store.CountAccounts(ctx)
	if err != nil {
		return nil, 0, err
	}

	return accounts, count, nil
}

// Update writes the new balance, then returns the authoritative account
// state so the handler never returns stale or partial data.
func (s *accountService) Update(ctx context.Context, params UpdateAccountParams) (db.Account, error) {
	if err := s.store.UpdateAccount(ctx, db.UpdateAccountParams{
		ID:      params.ID,
		Balance: params.Balance,
	}); err != nil {
		return db.Account{}, err
	}
	return s.store.GetAccount(ctx, params.ID)
}

func (s *accountService) Delete(ctx context.Context, id int64) error {
	return s.store.DeleteAccount(ctx, id)
}
