package basiq

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type PayRequestParams struct {
	PayRequests []struct {
		RequestID           string `json:"requestId"`
		Description         string `json:"description"`
		Amount              int    `json:"amount"`
		CollectFundsToFloat bool   `json:"collectFundsToFloat,omitempty"`
		CheckAccountBalance bool   `json:"checkAccountBalance,omitempty"`
		Payer               struct {
			PayerUserID         string `json:"payerUserId"`
			PayerBankBranchCode string `json:"payerBankBranchCode,omitempty"`
			PayerAccountNumber  string `json:"payerAccountNumber,omitempty"`
		} `json:"payer"`
	}
}

type PayRequestList struct {
	Type  string       `json:"type"`
	Count int          `json:"count"`
	Size  int          `json:"size"`
	Data  []PayRequest `json:"data"`
	Links struct {
		Self string `json:"self"`
		Next string `json:"next"`
	} `json:"links"`
}

type PayRequest struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	RequestID string `json:"requestId"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	Method    string `json:"method"`
	Status    string `json:"status"`
	Reason    struct {
		Code    string `json:"code"`
		Title   string `json:"title"`
		Details string `json:"details"`
	} `json:"reason"`
	Payer struct {
		PayerUserID         string `json:"payerUserId"`
		PayerAccountID      string `json:"payerAccountId"`
		PayerBankBranchCode string `json:"payerBankBranchCode"`
		PayerAccountNumber  string `json:"payerAccountNumber"`
	} `json:"payer"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Links       struct {
		Self string `json:"self"`
		Job  string `json:"job"`
	} `json:"links"`
}

type PayRequestJobList struct {
	Jobs []PayRequestJob `json:"jobs"`
}

type PayRequestJob struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	RequestID string `json:"requestId"`
	Links     struct {
		Self string `json:"self"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) PayRequest(ctx context.Context, payRequestID string) (PayRequest, error) {
	payRequest, err := a.payRequest(ctx, payRequestID)
	if err != nil && !IsUnauthorizedErr(err) {
		return payRequest, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return PayRequest{}, err
	}
	return a.payRequest(ctx, payRequestID)
}

func (a *API) PayRequests(ctx context.Context) (PayRequestList, error) {
	payRequestList, err := a.payRequests(ctx)
	if err != nil && !IsUnauthorizedErr(err) {
		return payRequestList, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return PayRequestList{}, err
	}
	return a.payRequests(ctx)
}

func (a *API) CreatePayRequest(ctx context.Context, params PayRequestParams) (PayRequestJobList, error) {
	payRequestJobList, err := a.createPayRequest(ctx, params)
	if err != nil && !IsUnauthorizedErr(err) {
		return payRequestJobList, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return PayRequestJobList{}, err
	}
	return a.createPayRequest(ctx, params)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) payRequest(ctx context.Context, payRequestID string) (PayRequest, error) {
	callURL, err := url.JoinPath(baseURL, "payments", "payrequests", payRequestID)
	if err != nil {
		return PayRequest{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return PayRequest{}, err
	}

	var payRequest PayRequest
	return payRequest, json.Unmarshal(data, &payRequest)
}

func (a *API) payRequests(ctx context.Context) (PayRequestList, error) {
	callURL, err := url.JoinPath(baseURL, "payments", "payrequests")
	if err != nil {
		return PayRequestList{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return PayRequestList{}, err
	}

	var list PayRequestList
	return list, json.Unmarshal(data, &list)
}

func (a *API) createPayRequest(ctx context.Context, params PayRequestParams) (PayRequestJobList, error) {
	callURL, err := url.JoinPath(baseURL, "payments", "payrequests")
	if err != nil {
		return PayRequestJobList{}, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return PayRequestJobList{}, err
	}

	data, err := a.makeCall(ctx, http.MethodPost, callURL, bytes.NewReader(payload))
	if err != nil {
		return PayRequestJobList{}, err
	}

	var list PayRequestJobList
	return list, json.Unmarshal(data, &list)
}
