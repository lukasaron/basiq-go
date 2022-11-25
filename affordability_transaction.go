package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type AffordabilityTransactionList struct {
	Type  string                     `json:"type"`
	Count int                        `json:"count"`
	Size  int                        `json:"size"`
	Data  []AffordabilityTransaction `json:"data"`
	Links struct {
		Self string `json:"self"`
		Next string `json:"next"`
	} `json:"links"`
}

type AffordabilityTransaction struct {
	Type            string `json:"type"`
	ID              string `json:"id"`
	Account         string `json:"account"`
	Amount          string `json:"amount"`
	Balance         string `json:"balance"`
	Class           string `json:"class"`
	Description     string `json:"description"`
	Direction       string `json:"direction"`
	Institution     string `json:"institution"`
	PostDate        string `json:"postDate"`
	Status          string `json:"status"`
	TransactionDate string `json:"transactionDate"`
	Links           struct {
		Account     string `json:"account"`
		Institution string `json:"institution"`
	} `json:"links"`
}

//---------------------------------------------------------------------------------------------------------------------

func (a *API) AffordabilityTransactions(ctx context.Context, userID, snapshotID string) (AffordabilityTransactionList, error) {
	affordabilityTransactionList, err := a.affordabilityTransactions(ctx, userID, snapshotID)
	if err != nil && !IsUnauthorizedErr(err) {
		return affordabilityTransactionList, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return AffordabilityTransactionList{}, err
	}
	return a.affordabilityTransactions(ctx, userID, snapshotID)
}

//---------------------------------------------------------------------------------------------------------------------

func (a *API) affordabilityTransactions(ctx context.Context, userID, snapshotID string) (AffordabilityTransactionList, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "affordability", snapshotID, "transactions")
	if err != nil {
		return AffordabilityTransactionList{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return AffordabilityTransactionList{}, err
	}

	var list AffordabilityTransactionList
	return list, json.Unmarshal(data, &list)
}
