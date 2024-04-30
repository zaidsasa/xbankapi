package types

import (
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	_ struct{} `type:"structure"`

	Name         string `json:"name"         validate:"minLen:3|maxLen:255"`
	Email        string `json:"email"        validate:"required|email|maxLen:255"`
	CurrencyCode string `json:"currencyCode" message:"only EUR currency code is currently supported" validate:"eq:EUR"`
}

type CreateAccountResponse struct {
	_ struct{} `type:"structure"`

	Account
}

type Account struct {
	_ struct{} `type:"structure"`

	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	CurrencyCode string    `json:"currencyCode"`
}

type AddMoneyRequest struct {
	_ struct{} `type:"structure"`

	Amount money.Amount `json:"amount" validate:"money_amount"`
}

type AddMoneyResponse struct {
	_ struct{} `type:"structure"`

	TransactionID uuid.UUID `json:"id"`
}

type TransferMoneyRequest struct {
	_ struct{} `type:"structure"`

	ReciverAccountID uuid.UUID    `json:"reciverAccountId" validate:"required"`
	Amount           money.Amount `json:"amount"           validate:"money_amount"`
}

type TransferMoneyResponse struct {
	_ struct{} `type:"structure"`

	TransactionID uuid.UUID `json:"id"`
}
