-- name: CreateAccount :one
INSERT INTO "account"(email, name, currency_code)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: AddTransaction :one
INSERT INTO "transaction"(account_id, amount, source_id)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: HasAccount :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            "account"
        WHERE
            account_id = $1);

-- name: GetAccount :one
SELECT
    *
FROM
    "account"
WHERE
    account_id = $1;

-- name: GetAccountTotalAmount :one
SELECT
    SUM(amount)::numeric
FROM
    "transaction"
WHERE
    account_id = $1;
