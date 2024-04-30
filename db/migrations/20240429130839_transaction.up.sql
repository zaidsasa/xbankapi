CREATE TABLE "transaction"(
    transaction_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id uuid NOT NULL REFERENCES account(account_id),
    amount numeric,
    source_id uuid REFERENCES TRANSACTION (transaction_id)
);
