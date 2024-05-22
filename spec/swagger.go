// Package classification xBankAPI
//
// Documentation of our xBank API.
//
//	 Schemes: http, https
//	 BasePath: /
//	 Version: 1.0.0
//	 Host: xbankapi.com
//
//	 Consumes:
//	 - application/json
//
//	 Produces:
//	 - application/json
//
//	 Security:
//	 - basic
//
//	SecurityDefinitions:
//	basic:
//	  type: basic
//
// swagger:meta
package spec

import "github.com/zaidsasa/xbankapi/internal/types"

// swagger:parameters CreateAccountRequest
type CreateAccountRequestBody struct {
	// in: path
	// required: true
	ID string `json:"id"`

	//  in: body
	// required: true
	Body types.CreateAccountRequest `json:"body"`
}

// swagger:response CreateAccountResponse
type CreateAccountResponse struct {
	//  in: body
	Body types.CreateAccountResponse `json:"body"`
}

// swagger:parameters AddMoneyRequest
type AddMoneyRequestBody struct {
	// in: path
	// required: true
	ID string `json:"id"`

	//  in: body
	// required: true
	Body types.AddMoneyRequest `json:"body"`
}

// swagger:response AddMoneyResponse
type AddMoneyResponseBody struct {
	//  in: body
	Body types.AddMoneyResponse `json:"body"`
}

// swagger:parameters TransferMoneyRequest
type TransferMoneyRequestBody struct {
	// in: path
	// required: true
	ID string `json:"id"`

	//  in: body
	// required: true
	Body types.TransferMoneyRequest `json:"body"`
}

// swagger:response TransferMoneyResponse
type TransferMoneyResponse struct {
	//  in: body
	Body types.TransferMoneyResponse `json:"body"`
}
