CREATE TABLE "account"(
    account_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email varchar(255) UNIQUE NOT NULL,
    name varchar(255) NOT NULL,
    currency_code varchar(3) NOT NULL
);
