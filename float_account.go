package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type FloatAccountList struct {
	Type  string         `json:"type"`
	Count int            `json:"count"`
	Size  int            `json:"size"`
	Data  []FloatAccount `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type FloatAccount struct {
	Type             string `json:"type"`
	ID               string `json:"id"`
	BankBranchCode   string `json:"bankBranchCode"`
	AccountNumber    string `json:"accountNumber"`
	AvailableBalance int    `json:"availableBalance"`
	Status           string `json:"status"`
	Links            struct {
		Self string `json:"self"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) FloatAccount(ctx context.Context, floatAccountID string) (FloatAccount, error) {
	floatAccount, err := a.floatAccount(ctx, floatAccountID)
	if err != nil && !IsUnauthorizedErr(err) {
		return floatAccount, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return FloatAccount{}, err
	}
	return a.floatAccount(ctx, floatAccountID)
}

func (a *API) FloatAccounts(ctx context.Context) ([]FloatAccount, error) {
	floatAccounts, err := a.floatAccounts(ctx)
	if err != nil && !IsUnauthorizedErr(err) {
		return floatAccounts, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return nil, err
	}
	return a.floatAccounts(ctx)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) floatAccount(ctx context.Context, floatAccountID string) (FloatAccount, error) {
	callURL, err := url.JoinPath(baseURL, "payments", "float-accounts", floatAccountID)
	if err != nil {
		return FloatAccount{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return FloatAccount{}, err
	}

	var floatAccount FloatAccount
	return floatAccount, json.Unmarshal(data, &floatAccount)
}

func (a *API) floatAccounts(ctx context.Context) ([]FloatAccount, error) {
	callURL, err := url.JoinPath(baseURL, "payments", "float-accounts")
	if err != nil {
		return nil, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return nil, err
	}

	var list FloatAccountList
	if err = json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list.Data, nil
}
