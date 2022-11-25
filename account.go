package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type AccountList struct {
	Type  string    `json:"type"`
	Data  []Account `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type Account struct {
	Type           string `json:"type"`
	ID             string `json:"id"`
	AccountHolder  string `json:"accountHolder"`
	AccountNo      string `json:"accountNo"`
	AvailableFunds string `json:"availableFunds"`
	Balance        string `json:"balance"`
	Class          []struct {
		Type    string `json:"type"`
		Product string `json:"product"`
	} `json:"class"`
	Connection           string `json:"connection"`
	Currency             string `json:"currency"`
	Institution          string `json:"institution"`
	LastUpdated          string `json:"lastUpdated"`
	Name                 string `json:"name"`
	Status               string `json:"status"`
	TransactionIntervals []struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"transactionIntervals"`
	Links struct {
		Institution  string `json:"institution"`
		Transactions string `json:"transactions"`
		Self         string `json:"self"`
	} `json:"links"`
}

//---------------------------------------------------------------------------------------------------------------------

func (c *Client) Account(ctx context.Context, userID, accountID string) (Account, error) {
	account, err := c.account(ctx, userID, accountID)
	if err != nil && !IsUnauthorizedErr(err) {
		return account, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return Account{}, err
	}
	return c.account(ctx, userID, accountID)
}

func (c *Client) Accounts(ctx context.Context, userID string) (AccountList, error) {
	accountList, err := c.accounts(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return accountList, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return AccountList{}, err
	}
	return c.accounts(ctx, userID)
}

//---------------------------------------------------------------------------------------------------------------------

func (c *Client) account(ctx context.Context, userID, accountID string) (Account, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "accounts", accountID)
	if err != nil {
		return Account{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return Account{}, err
	}

	var account Account
	return account, json.Unmarshal(data, &account)
}

func (c *Client) accounts(ctx context.Context, userID string) (AccountList, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "accounts")
	if err != nil {
		return AccountList{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return AccountList{}, err
	}

	var list AccountList
	return list, json.Unmarshal(data, &list)
}
