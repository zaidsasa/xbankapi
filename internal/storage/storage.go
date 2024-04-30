package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBConnection interface {
	Ping(ctx context.Context) error
	Begin(ctx context.Context) (pgx.Tx, error)
}

type AccountStore interface {
	AddTransaction(ctx context.Context, arg AddTransactionParams) (Transaction, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	GetAccount(ctx context.Context, accountID uuid.UUID) (Account, error)
	GetAccountTotalAmount(ctx context.Context, accountID uuid.UUID) (pgtype.Numeric, error)
	HasAccount(ctx context.Context, accountID uuid.UUID) (bool, error)
}

var AccountStoreWithTx = func(tx pgx.Tx) AccountStore {
	return &Queries{
		db: tx,
	}
}
