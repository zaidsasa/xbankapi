// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const addTransaction = `-- name: AddTransaction :one
INSERT INTO "transaction"(account_id, amount, source_id)
    VALUES ($1, $2, $3)
RETURNING
    transaction_id, account_id, amount, source_id
`

type AddTransactionParams struct {
	AccountID uuid.UUID
	Amount    pgtype.Numeric
	SourceID  uuid.NullUUID
}

func (q *Queries) AddTransaction(ctx context.Context, arg AddTransactionParams) (Transaction, error) {
	row := q.db.QueryRow(ctx, addTransaction, arg.AccountID, arg.Amount, arg.SourceID)
	var i Transaction
	err := row.Scan(
		&i.TransactionID,
		&i.AccountID,
		&i.Amount,
		&i.SourceID,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO "account"(email, name, currency_code)
    VALUES ($1, $2, $3)
RETURNING
    account_id, email, name, currency_code
`

type CreateAccountParams struct {
	Email        string
	Name         string
	CurrencyCode string
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRow(ctx, createAccount, arg.Email, arg.Name, arg.CurrencyCode)
	var i Account
	err := row.Scan(
		&i.AccountID,
		&i.Email,
		&i.Name,
		&i.CurrencyCode,
	)
	return i, err
}

const getAccount = `-- name: GetAccount :one
SELECT
    account_id, email, name, currency_code
FROM
    "account"
WHERE
    account_id = $1
`

func (q *Queries) GetAccount(ctx context.Context, accountID uuid.UUID) (Account, error) {
	row := q.db.QueryRow(ctx, getAccount, accountID)
	var i Account
	err := row.Scan(
		&i.AccountID,
		&i.Email,
		&i.Name,
		&i.CurrencyCode,
	)
	return i, err
}

const getAccountTotalAmount = `-- name: GetAccountTotalAmount :one
SELECT
    SUM(amount)::numeric
FROM
    "transaction"
WHERE
    account_id = $1
`

func (q *Queries) GetAccountTotalAmount(ctx context.Context, accountID uuid.UUID) (pgtype.Numeric, error) {
	row := q.db.QueryRow(ctx, getAccountTotalAmount, accountID)
	var column_1 pgtype.Numeric
	err := row.Scan(&column_1)
	return column_1, err
}

const hasAccount = `-- name: HasAccount :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            "account"
        WHERE
            account_id = $1)
`

func (q *Queries) HasAccount(ctx context.Context, accountID uuid.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, hasAccount, accountID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
