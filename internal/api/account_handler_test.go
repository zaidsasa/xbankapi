package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zaidsasa/xbankapi/internal/api/mocks"
	"github.com/zaidsasa/xbankapi/internal/types"
	"github.com/zaidsasa/xbankapi/internal/validator"
)

func TestNewAccountHandler(t *testing.T) {
	t.Parallel()

	got := NewAccountHandler(&ImplAccountService{})
	assert.NotNil(t, got)
}

func TestAccountHandler_createAccount(t *testing.T) {
	t.Parallel()

	type args struct {
		body types.CreateAccountRequest
	}

	tests := []struct {
		name           string
		args           args
		mock           func(*mocks.MockAccountService, args)
		wantStatusCode int
		want           string
	}{
		{
			name: "failed when email is invalid",
			args: args{
				body: types.CreateAccountRequest{
					Name:         "name",
					Email:        "wrong",
					CurrencyCode: "EUR",
				},
			},
			wantStatusCode: http.StatusBadRequest,
			want:           `{"email":{"email":"email value is an invalid email address"}}`,
		},
		{
			name: "failed when currency code is invalid",
			args: args{
				body: types.CreateAccountRequest{
					Name:         "name",
					Email:        "test@mail.com",
					CurrencyCode: "invalid currency",
				},
			},
			wantStatusCode: http.StatusBadRequest,
			want:           `{"currencyCode":{"eq":"only EUR currency code is currently supported"}}`,
		},
		{
			name: "failed when name is invalid",
			args: args{
				body: types.CreateAccountRequest{
					Name:         "n",
					Email:        "test@mail.com",
					CurrencyCode: "EUR",
				},
			},
			wantStatusCode: http.StatusBadRequest,
			want:           `{"name":{"minLen":"name min length is 3"}}`,
		},
		{
			name: "failed when create account returns an internal error",
			args: args{
				body: types.CreateAccountRequest{
					Name:         "name",
					Email:        "test@mail.com",
					CurrencyCode: "EUR",
				},
			},

			mock: func(mas *mocks.MockAccountService, _ args) {
				mas.EXPECT().CreateAccount(mock.Anything, mock.Anything).
					Return(types.CreateAccountResponse{}, ErrInternal).Once()
			},
			wantStatusCode: http.StatusInternalServerError,
			want: `{"message":"internal server error"}
`,
		},
		{
			name: "success when creating an account",
			args: args{
				body: types.CreateAccountRequest{
					Name:         "name",
					Email:        "test@mail.com",
					CurrencyCode: "EUR",
				},
			},

			mock: func(mas *mocks.MockAccountService, a args) {
				mas.EXPECT().CreateAccount(mock.Anything, mock.Anything).Return(types.CreateAccountResponse{
					Account: types.Account{
						ID:           wantAccountID,
						Name:         a.body.Name,
						Email:        a.body.Email,
						CurrencyCode: a.body.CurrencyCode,
					},
				}, nil).Once()
			},
			wantStatusCode: http.StatusOK,
			want: `{"id":"12345678-1234-1234-1234-123456789001","name":"name","email":"test@mail.com","currencyCode":"EUR"}
`,
		},
	}

	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			body, err := json.Marshal(tt.args.body)
			assert.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
			w := httptest.NewRecorder()

			accountServiceMock := mocks.NewMockAccountService(t)

			if tt.mock != nil {
				tt.mock(accountServiceMock, tt.args)
			}

			accountHandler := NewAccountHandler(accountServiceMock)
			accountHandler.createAccount(w, r)

			res := w.Result()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			defer res.Body.Close()

			got, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestAccountHandler_addMoney(t *testing.T) {
	validator.ConfigureDefaultValidator()

	t.Parallel()

	type args struct {
		accountID uuid.UUID
		body      types.AddMoneyRequest
	}

	tests := []struct {
		name           string
		args           args
		mock           func(*mocks.MockAccountService, args)
		wantStatusCode int
		want           string
	}{
		{
			name: "failed when amount is invalid",
			args: args{
				accountID: wantAccountID,
				body: types.AddMoneyRequest{
					Amount: -1,
				},
			},
			wantStatusCode: http.StatusBadRequest,
			want:           `{"amount":{"money_amount":"amount field did not pass validation"}}`,
		},
		{
			name: "success when transaction is created",
			args: args{
				accountID: wantAccountID,
				body: types.AddMoneyRequest{
					Amount: 111,
				},
			},
			mock: func(mas *mocks.MockAccountService, _ args) {
				mas.EXPECT().AddMoney(mock.Anything, mock.Anything, wantAccountID).Return(types.AddMoneyResponse{
					TransactionID: wantTrnasactionID,
				}, nil).Once()
			},
			wantStatusCode: http.StatusOK,
			want: `{"id":"12345678-1234-1234-1234-123456789002"}
`,
		},
	}
	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			body, err := json.Marshal(tt.args.body)
			assert.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, "/accounts/:d/transactions", bytes.NewReader(body))
			r.SetPathValue(pathValueID, tt.args.accountID.String())

			w := httptest.NewRecorder()

			accountServiceMock := mocks.NewMockAccountService(t)

			if tt.mock != nil {
				tt.mock(accountServiceMock, tt.args)
			}

			accountHandler := NewAccountHandler(accountServiceMock)
			accountHandler.addMoney(w, r)

			res := w.Result()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			defer res.Body.Close()

			got, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestAccountHandler_transferMoney(t *testing.T) {
	validator.ConfigureDefaultValidator()

	t.Parallel()

	type args struct {
		accountID uuid.UUID
		body      types.TransferMoneyRequest
	}

	tests := []struct {
		name           string
		args           args
		mock           func(*mocks.MockAccountService)
		wantStatusCode int
		want           string
	}{
		{
			name: "failed when amount is invalid",
			args: args{
				accountID: wantAccountID,
				body: types.TransferMoneyRequest{
					Amount: -1,
				},
			},
			wantStatusCode: http.StatusBadRequest,
			want:           `{"amount":{"money_amount":"amount field did not pass validation"}}`,
		},
		{
			name: "success when money is transferred",
			args: args{
				accountID: wantAccountID,
				body: types.TransferMoneyRequest{
					Amount: 100,
				},
			},
			mock: func(mas *mocks.MockAccountService) {
				mas.EXPECT().TransferMoney(mock.Anything, mock.Anything, wantAccountID).
					Return(types.TransferMoneyResponse{
						TransactionID: wantReciverTransactionID,
					}, nil).Once()
			},
			wantStatusCode: http.StatusOK,
			want: `{"id":"12345678-1234-1234-1234-123456789004"}
`,
		},
	}
	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			body, err := json.Marshal(tt.args.body)
			assert.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, "/accounts/:d/transactions/transfer", bytes.NewReader(body))
			r.SetPathValue(pathValueID, tt.args.accountID.String())

			w := httptest.NewRecorder()

			accountServiceMock := mocks.NewMockAccountService(t)

			if tt.mock != nil {
				tt.mock(accountServiceMock)
			}

			accountHandler := NewAccountHandler(accountServiceMock)
			accountHandler.transferMoney(w, r)

			res := w.Result()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			defer res.Body.Close()

			got, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}
