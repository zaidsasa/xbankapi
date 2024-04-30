package api

import (
	"context"
	"errors"
	"log/slog"
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zaidsasa/xbankapi/internal/storage"
	storageMocks "github.com/zaidsasa/xbankapi/internal/storage/mocks"
	"github.com/zaidsasa/xbankapi/internal/types"
	txMocks "github.com/zaidsasa/xbankapi/mocks/github.com/jackc/pgx/v5"
)

var (
	wantAccountID            = uuid.MustParse("12345678-1234-1234-1234-123456789001")
	wantTrnasactionID        = uuid.MustParse("12345678-1234-1234-1234-123456789002")
	wantReciverAccountID     = uuid.MustParse("12345678-1234-1234-1234-123456789003")
	wantReciverTransactionID = uuid.MustParse("12345678-1234-1234-1234-123456789004")
)

func TestNewAccountService(t *testing.T) {
	t.Parallel()

	got := NewAccountService(&pgxpool.Pool{}, storageMocks.NewMockAccountStore(t), slog.Default())
	assert.NotNil(t, got)
}

func TestAccountService_CreateAccount(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *types.CreateAccountRequest
	}

	tests := []struct {
		name    string
		args    args
		mock    func(*storageMocks.MockAccountStore, args)
		want    types.CreateAccountResponse
		wantErr error
	}{
		{
			name: "failed when creating an account returns an error",
			args: args{
				ctx: context.Background(),
				req: &types.CreateAccountRequest{},
			},
			mock: func(accountStorageMock *storageMocks.MockAccountStore, a args) {
				accountStorageMock.EXPECT().CreateAccount(a.ctx, mock.Anything).
					Return(storage.Account{}, errors.New("error")).Once()
			},
			wantErr: ErrInternal,
		},
		{
			name: "success when creating an account",
			args: args{
				ctx: context.Background(),
				req: &types.CreateAccountRequest{
					Name:         "test",
					Email:        "test@mail.com",
					CurrencyCode: "EUR",
				},
			},
			mock: func(accountStorageMock *storageMocks.MockAccountStore, a args) {
				accountStorageMock.EXPECT().CreateAccount(a.ctx, mock.Anything).Return(storage.Account{
					AccountID:    wantAccountID,
					Name:         a.req.Name,
					Email:        a.req.Email,
					CurrencyCode: a.req.CurrencyCode,
				}, nil).Once()
			},
			want: types.CreateAccountResponse{
				Account: types.Account{
					ID:           wantAccountID,
					Name:         "test",
					Email:        "test@mail.com",
					CurrencyCode: "EUR",
				},
			},
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			accountStorageMock := storageMocks.NewMockAccountStore(t)
			connMock := storageMocks.NewMockDBConnection(t)
			logger := slog.Default()

			accountService := NewAccountService(connMock, accountStorageMock, logger)

			tt.mock(accountStorageMock, tt.args)
			got, err := accountService.CreateAccount(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.want, got)

			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccountService_AddMoney(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		req       *types.AddMoneyRequest
		accountID uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		mock    func(*storageMocks.MockAccountStore, args)
		want    types.AddMoneyResponse
		wantErr error
	}{
		{
			name: "failed when account not found",
			args: args{
				ctx: context.Background(),
				req: &types.AddMoneyRequest{
					Amount: 100,
				},
				accountID: uuid.New(),
			},
			mock: func(accountStorageMock *storageMocks.MockAccountStore, a args) {
				accountStorageMock.EXPECT().HasAccount(a.ctx, a.accountID).
					Return(false, nil).Once()
			},
			wantErr: ErrAccountNotFound,
		},
		{
			name: "failed when create transaction returns an error",
			args: args{
				ctx: context.Background(),
				req: &types.AddMoneyRequest{
					Amount: 100,
				},
				accountID: uuid.New(),
			},
			mock: func(accountStorageMock *storageMocks.MockAccountStore, args args) {
				accountStorageMock.EXPECT().HasAccount(args.ctx, args.accountID).
					Return(true, nil).Once()
				accountStorageMock.EXPECT().AddTransaction(args.ctx, mock.Anything).
					Return(storage.Transaction{}, errors.New("error")).Once()
			},
			wantErr: ErrInternal,
		},
		{
			name: "success when account exists",
			args: args{
				ctx: context.Background(),
				req: &types.AddMoneyRequest{
					Amount: 100,
				},
				accountID: uuid.New(),
			},
			mock: func(accountStorageMock *storageMocks.MockAccountStore, args args) {
				accountStorageMock.EXPECT().HasAccount(args.ctx, args.accountID).
					Return(true, nil).Once()
				accountStorageMock.EXPECT().AddTransaction(args.ctx, mock.Anything).
					Return(storage.Transaction{
						TransactionID: wantTrnasactionID,
						AccountID:     args.accountID,
						Amount:        pgtype.Numeric{Int: big.NewInt(args.req.Amount), Exp: -2, Valid: true},
					}, nil).Once()
			},
			want: types.AddMoneyResponse{
				TransactionID: wantTrnasactionID,
			},
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			accountStorageMock := storageMocks.NewMockAccountStore(t)
			connMock := storageMocks.NewMockDBConnection(t)
			logger := slog.Default()

			tt.mock(accountStorageMock, tt.args)

			accountService := NewAccountService(connMock, accountStorageMock, logger)
			got, err := accountService.AddMoney(tt.args.ctx, tt.args.req, tt.args.accountID)

			assert.Equal(t, tt.want, got)

			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccountService_TransferMoney(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		req       *types.TransferMoneyRequest
		accountID uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		mock    func(*storageMocks.MockAccountStore, *storageMocks.MockDBConnection, args)
		want    types.TransferMoneyResponse
		wantErr error
	}{
		{
			name: "failed when get account returns an error",
			args: args{
				req: &types.TransferMoneyRequest{
					ReciverAccountID: wantReciverAccountID,
					Amount:           200,
				},
				accountID: wantAccountID,
			},
			mock: func(accountStorageMock *storageMocks.MockAccountStore, _ *storageMocks.MockDBConnection, a args) {
				accountStorageMock.EXPECT().GetAccount(a.ctx, a.accountID).
					Return(storage.Account{}, errors.New("error"))
			},
			wantErr: ErrInternal,
		},
		{
			name: "failed when get account total amount returns an error",
			args: args{
				req: &types.TransferMoneyRequest{
					ReciverAccountID: wantReciverAccountID,
					Amount:           200,
				},
				accountID: wantAccountID,
			},
			mock: func(accountStorageMock *storageMocks.MockAccountStore, _ *storageMocks.MockDBConnection, a args) {
				accountStorageMock.EXPECT().GetAccount(a.ctx, a.accountID).
					Return(storage.Account{}, nil).Once()
				accountStorageMock.EXPECT().GetAccountTotalAmount(a.ctx, a.accountID).
					Return(pgtype.Numeric{}, errors.New("error")).Once()
			},
			wantErr: ErrInternal,
		},
		{
			name: "failed when insufficient account balance",
			args: args{
				req: &types.TransferMoneyRequest{
					ReciverAccountID: wantReciverAccountID,
					Amount:           200,
				},
				accountID: wantAccountID,
			},
			mock: func(accountStorageMock *storageMocks.MockAccountStore, _ *storageMocks.MockDBConnection, a args) {
				accountStorageMock.EXPECT().GetAccount(a.ctx, a.accountID).
					Return(storage.Account{}, nil).Once()
				accountStorageMock.EXPECT().GetAccountTotalAmount(a.ctx, a.accountID).
					Return(pgtype.Numeric{Int: big.NewInt(200), Exp: -2}, nil).Once()
			},
			wantErr: ErrInsufficientAccountBalance,
		},
		{
			name: "success when money transfer is succeeded",
			args: args{
				ctx: context.Background(),
				req: &types.TransferMoneyRequest{
					ReciverAccountID: wantReciverAccountID,
					Amount:           200,
				},
				accountID: wantAccountID,
			},
			mock: func(accountStorageMock *storageMocks.MockAccountStore, conn *storageMocks.MockDBConnection, a args) {
				accountStorageMock.EXPECT().GetAccount(a.ctx, a.accountID).
					Return(storage.Account{}, nil).Once()
				accountStorageMock.EXPECT().GetAccountTotalAmount(a.ctx, a.accountID).
					Return(pgtype.Numeric{Int: big.NewInt(201), Exp: -2}, nil).Once()

				tx := txMocks.NewMockTx(t)
				conn.EXPECT().Begin(a.ctx).Return(tx, nil).Once()

				storage.AccountStoreWithTx = func(_ pgx.Tx) storage.AccountStore {
					return accountStorageMock
				}

				accountStorageMock.EXPECT().AddTransaction(a.ctx, mock.Anything).Return(storage.Transaction{}, nil).Once()

				accountStorageMock.EXPECT().AddTransaction(a.ctx, mock.Anything).Return(storage.Transaction{
					TransactionID: wantReciverTransactionID,
				}, nil).Once()

				tx.EXPECT().Commit(a.ctx).Return(nil).Once()
			},
			want: types.TransferMoneyResponse{
				TransactionID: wantReciverTransactionID,
			},
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			accountStorageMock := storageMocks.NewMockAccountStore(t)
			connMock := storageMocks.NewMockDBConnection(t)
			logger := slog.Default()

			tt.mock(accountStorageMock, connMock, tt.args)

			accountService := NewAccountService(connMock, accountStorageMock, logger)
			got, err := accountService.TransferMoney(tt.args.ctx, tt.args.req, tt.args.accountID)
			assert.Equal(t, tt.want, got)

			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
