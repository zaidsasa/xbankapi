package api

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zaidsasa/xbankapi/internal/logger"
	"github.com/zaidsasa/xbankapi/internal/storage"
	"github.com/zaidsasa/xbankapi/internal/types"
)

const (
	pqErrorForeignKeyViolation = "23503"
	pqErrorAlreadyExist        = "23505"
)

var (
	ErrInsufficientAccountBalance = errors.New("insufficient account balance")
	ErrAccountNotFound            = errors.New("account not found")
	ErrRecieverAccountNotFound    = errors.New("reciver account not found")
	ErrInternal                   = errors.New("internal error")
	ErrAccountAlreadyExist        = errors.New("account already exists")
)

type AccountService interface {
	CreateAccount(ctx context.Context, req *types.CreateAccountRequest) (types.CreateAccountResponse, error)
	AddMoney(ctx context.Context, req *types.AddMoneyRequest, accountID uuid.UUID) (types.AddMoneyResponse, error)
	TransferMoney(
		ctx context.Context, req *types.TransferMoneyRequest, accountID uuid.UUID) (types.TransferMoneyResponse, error)
}

type ImplAccountService struct {
	logger logger.Logger
	xlock  map[uuid.UUID]sync.Locker
	conn   storage.DBConnection
	store  storage.AccountStore
}

// NewAccountService returns a new ImplAccountService.
func NewAccountService(
	conn storage.DBConnection,
	store storage.AccountStore,
	logger logger.Logger,
) *ImplAccountService {
	return &ImplAccountService{
		logger: logger,
		xlock:  make(map[uuid.UUID]sync.Locker),
		conn:   conn,
		store:  store,
	}
}

// CreateAccount creates a bank account.
// returns CreateAccountResponse.
func (a *ImplAccountService) CreateAccount(
	ctx context.Context,
	req *types.CreateAccountRequest,
) (types.CreateAccountResponse, error) {
	account, err := a.store.CreateAccount(ctx, storage.CreateAccountParams{
		Email:        req.Email,
		Name:         req.Name,
		CurrencyCode: req.CurrencyCode,
	})
	if err != nil {
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr); pgErr.Code == pqErrorAlreadyExist {
			return types.CreateAccountResponse{}, ErrAccountAlreadyExist
		}

		a.logger.Error("failed to create account", "error", err)

		return types.CreateAccountResponse{}, ErrInternal
	}

	return types.CreateAccountResponse{
		Account: types.Account{
			ID:           account.AccountID,
			Name:         account.Name,
			Email:        account.Email,
			CurrencyCode: req.CurrencyCode,
		},
	}, nil
}

// AddMoney add money to bank account.
// returns AddMoneyResponse.
func (a *ImplAccountService) AddMoney(
	ctx context.Context,
	req *types.AddMoneyRequest,
	accountID uuid.UUID,
) (types.AddMoneyResponse, error) {
	if err := a.hasAccount(ctx, accountID); err != nil {
		return types.AddMoneyResponse{}, err
	}

	t, err := a.store.AddTransaction(ctx, storage.AddTransactionParams{
		AccountID: accountID,
		Amount:    pgtype.Numeric{Int: big.NewInt(req.Amount), Exp: -2, Valid: true},
	})
	if err != nil {
		a.logger.Error("failed to add money", "error", err)

		return types.AddMoneyResponse{}, ErrInternal
	}

	return types.AddMoneyResponse{
		TransactionID: t.TransactionID,
	}, nil
}

// TransferMoney transfers money from a bank account to another.
// returns TransferMoneyResponse.
func (a *ImplAccountService) TransferMoney(
	ctx context.Context,
	req *types.TransferMoneyRequest,
	accountID uuid.UUID,
) (types.TransferMoneyResponse, error) {
	account, err := a.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.TransferMoneyResponse{}, ErrAccountNotFound
		}

		a.logger.Error("failed to fetch account", "error", err)

		return types.TransferMoneyResponse{}, ErrInternal
	}

	a.lock(accountID)
	defer a.unlock(accountID)

	totalAmount, err := a.store.GetAccountTotalAmount(ctx, accountID)
	if err != nil {
		a.logger.Error("failed to get account total amount", "error", err)

		return types.TransferMoneyResponse{}, ErrInternal
	}

	if err = validateTotalBalanceForMoneyTransfer(
		totalAmount,
		req.Amount,
		account.CurrencyCode); err != nil {
		a.logger.Error("failed to calculate expected total balance", "error", err)

		return types.TransferMoneyResponse{}, err
	}

	tx, err := a.conn.Begin(ctx)
	if err != nil {
		a.logger.Error("failed to begin transaction", "error", err)

		return types.TransferMoneyResponse{}, ErrInternal
	}

	s := storage.AccountStoreWithTx(tx)

	t, err := s.AddTransaction(ctx, storage.AddTransactionParams{
		AccountID: accountID, Amount: pgtype.Numeric{Int: big.NewInt(req.Amount * -1), Exp: -2, Valid: true},
	})
	if err != nil {
		a.logger.Error("failed to add transaction", "error", err)

		return types.TransferMoneyResponse{}, ErrInternal
	}

	reciverTransaction, err := s.AddTransaction(ctx, storage.AddTransactionParams{
		AccountID: req.ReciverAccountID,
		Amount:    pgtype.Numeric{Int: big.NewInt(req.Amount), Exp: -2, Valid: true},
		SourceID:  uuid.NullUUID{UUID: t.TransactionID, Valid: true},
	})
	if err != nil {
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr); pgErr.Code == pqErrorForeignKeyViolation {
			return types.TransferMoneyResponse{}, ErrRecieverAccountNotFound
		}

		a.logger.Error("failed to add transaction", "error", err)

		return types.TransferMoneyResponse{}, ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		a.logger.Error("failed to commit transaction", "error", err)

		return types.TransferMoneyResponse{}, ErrInternal
	}

	return types.TransferMoneyResponse{TransactionID: reciverTransaction.TransactionID}, nil
}

func validateTotalBalanceForMoneyTransfer(
	totalAmount pgtype.Numeric,
	transferableAmount int64,
	currencyCode string,
) error {
	if totalAmount.Int == nil {
		return ErrInsufficientAccountBalance
	}

	totalMoney := money.New(totalAmount.Int.Int64(), currencyCode)
	transferAmountMoney := money.New(transferableAmount, currencyCode)

	res, err := totalMoney.Subtract(transferAmountMoney)
	if err != nil {
		return fmt.Errorf("failed to subtract amount: %w", err)
	}

	if !res.IsPositive() {
		return ErrInsufficientAccountBalance
	}

	return nil
}

func (a *ImplAccountService) hasAccount(ctx context.Context, accountID uuid.UUID) error {
	ok, err := a.store.HasAccount(ctx, accountID)
	if err != nil {
		a.logger.Error("failed to check account", "error", err)

		return ErrInternal
	}

	if !ok {
		return ErrAccountNotFound
	}

	return nil
}

// BUG: Current lock implementation will lead to a memory leakage and doesn't support high availability.
func (a *ImplAccountService) getLock(accountID uuid.UUID) sync.Locker {
	if a.xlock[accountID] == nil {
		a.xlock[accountID] = &sync.Mutex{}
	}

	return a.xlock[accountID]
}

func (a *ImplAccountService) lock(accountID uuid.UUID) {
	lock := a.getLock(accountID)

	lock.Lock()
}

func (a *ImplAccountService) unlock(accountID uuid.UUID) {
	lock := a.getLock(accountID)

	lock.Unlock()
}
