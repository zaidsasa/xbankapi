package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gookit/validate"
	"github.com/zaidsasa/xbankapi/internal/types"
)

const (
	createAccountRoute = "POST /accounts"
	addMoneyRoute      = "POST /accounts/{id}/transactions"
	transferMoneyRoute = "POST /accounts/{id}/transactions/transfer"

	pathValueID = "id"
)

type AccountHandler struct {
	service AccountService
}

// NewAccountHandler returns a new AccountHandler.
func NewAccountHandler(service AccountService) *AccountHandler {
	return &AccountHandler{
		service: service,
	}
}

// Register routes.
func (h *AccountHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(createAccountRoute, h.createAccount)
	mux.HandleFunc(addMoneyRoute, h.addMoney)
	mux.HandleFunc(transferMoneyRoute, h.transferMoney)
}

func (h *AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &types.CreateAccountRequest{}

	if err := h.decode(r, req); err != nil {
		handleError(w, err, http.StatusBadRequest)

		return
	}

	if v := validate.Struct(req); !v.Validate() {
		handleError(w, v.Errors, http.StatusBadRequest)

		return
	}

	res, err := h.service.CreateAccount(ctx, req)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)

		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
}

func (h *AccountHandler) addMoney(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &types.AddMoneyRequest{}
	if err := h.decode(r, req); err != nil {
		handleError(w, err, http.StatusBadRequest)

		return
	}

	if v := validate.Struct(req); !v.Validate() {
		handleError(w, v.Errors, http.StatusBadRequest)

		return
	}

	accountID, err := uuid.Parse(r.PathValue(pathValueID))
	if err != nil {
		handleError(w, err, http.StatusBadRequest)

		return
	}

	res, err := h.service.AddMoney(ctx, req, accountID)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)

		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
}

func (h *AccountHandler) transferMoney(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &types.TransferMoneyRequest{}
	if err := h.decode(r, req); err != nil {
		handleError(w, err, http.StatusBadRequest)

		return
	}

	if v := validate.Struct(req); !v.Validate() {
		handleError(w, v.Errors, http.StatusBadRequest)

		return
	}

	accountID, err := uuid.Parse(r.PathValue(pathValueID))
	if err != nil {
		handleError(w, err, http.StatusBadRequest)

		return
	}

	res, err := h.service.TransferMoney(ctx, req, accountID)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)

		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
}

func (h *AccountHandler) decode(req *http.Request, obj any) error {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&obj)
	if err != nil {
		return fmt.Errorf("unable to decode: %w", err)
	}

	return nil
}

type jsonError struct {
	Messsage string `json:"message"`
}

// TODO: Create a proper error types containing cause and statusCode.
func handleError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	violationError := &validate.Errors{}
	if errors.As(err, violationError) {
		w.WriteHeader(code)
		_, _ = w.Write(violationError.JSON())

		return
	}

	var jsonErr jsonError

	if errors.Is(err, ErrInternal) {
		code = http.StatusInternalServerError

		jsonErr = jsonError{
			Messsage: "internal server error",
		}
	} else {
		jsonErr = jsonError{
			Messsage: err.Error(),
		}
	}

	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(jsonErr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
