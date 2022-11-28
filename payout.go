package basiq

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type PayoutParams struct {
	RequestID   string `json:"requestId"`
	Method      string `json:"method,omitempty"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Payee       struct {
		PayeeUserID         string `json:"payeeUserId"`
		PayeeBankBranchCode string `json:"payeeBankBranchCode"`
		PayeeAccountNumber  string `json:"payeeAccountNumber"`
	} `json:"payee"`
}

type PayoutList struct {
	Type  string   `json:"type"`
	Count int      `json:"count"`
	Size  int      `json:"size"`
	Data  []Payout `json:"data"`
	Links struct {
		Self string `json:"self"`
		Next string `json:"next"`
	} `json:"links"`
}

type Payout struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	RequestID string `json:"requestId"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	Method    string `json:"method"`
	Status    string `json:"status"`
	Reason    struct {
		Code   string `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"reason"`
	Payee struct {
		PayeeUserID         string `json:"payeeUserId"`
		PayeeAccountID      string `json:"payeeAccountId"`
		PayeeBankBranchCode string `json:"payeeBankBranchCode"`
		PayeeAccountNumber  string `json:"payeeAccountNumber"`
	} `json:"payee"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Links       struct {
		Self string `json:"self"`
		Job  string `json:"job"`
	} `json:"links"`
}

type PayoutJobList struct {
	Jobs []PayoutJob `json:"jobs"`
}

type PayoutJob struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	RequestID string `json:"requestId"`
	Links     struct {
		Self string `json:"self"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) Payout(ctx context.Context, payoutID string) (Payout, error) {
	payout, err := a.payout(ctx, payoutID)
	if err != nil && !IsUnauthorizedErr(err) {
		return payout, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return Payout{}, err
	}
	return a.payout(ctx, payoutID)
}

func (a *API) Payouts(ctx context.Context) ([]Payout, error) {
	payouts, err := a.payouts(ctx)
	if err != nil && !IsUnauthorizedErr(err) {
		return payouts, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return nil, err
	}
	return a.payouts(ctx)
}

func (a *API) CreatePayout(ctx context.Context, params PayoutParams) ([]PayoutJob, error) {
	payoutJobs, err := a.createPayout(ctx, params)
	if err != nil && !IsUnauthorizedErr(err) {
		return payoutJobs, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return nil, err
	}
	return a.createPayout(ctx, params)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) payout(ctx context.Context, payoutID string) (Payout, error) {
	callURL, err := url.JoinPath(baseURL, "payments", "payouts", payoutID)
	if err != nil {
		return Payout{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return Payout{}, err
	}

	var payout Payout
	return payout, json.Unmarshal(data, &payout)
}

func (a *API) payouts(ctx context.Context) ([]Payout, error) {
	callURL, err := url.JoinPath(baseURL, "payments", "payouts")
	if err != nil {
		return nil, err
	}

	var payouts []Payout
	for callURL != "" {
		data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
		if err != nil {
			return nil, err
		}
		var list PayoutList
		if err = json.Unmarshal(data, &list); err != nil {
			return nil, err
		}
		if len(list.Data) > 0 {
			payouts = append(payouts, list.Data...)
		}
		callURL = list.Links.Next
	}

	return payouts, err
}

func (a *API) createPayout(ctx context.Context, params PayoutParams) ([]PayoutJob, error) {
	callURL, err := url.JoinPath(baseURL, "payments", "payouts")
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	data, err := a.makeCall(ctx, http.MethodPost, callURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	var list PayoutJobList
	if err = json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list.Jobs, nil
}
