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

func (a *API) AffordabilityTransactions(ctx context.Context, userID, snapshotID string) ([]AffordabilityTransaction, error) {
	affordabilityTransactions, err := a.affordabilityTransactions(ctx, userID, snapshotID)
	if err != nil && !IsUnauthorizedErr(err) {
		return affordabilityTransactions, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return nil, err
	}
	return a.affordabilityTransactions(ctx, userID, snapshotID)
}

//---------------------------------------------------------------------------------------------------------------------

func (a *API) affordabilityTransactions(ctx context.Context, userID, snapshotID string) ([]AffordabilityTransaction, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "affordability", snapshotID, "transactions")
	if err != nil {
		return nil, err
	}

	var transactions []AffordabilityTransaction
	for callURL != "" {
		data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
		if err != nil {
			return nil, err
		}
		var list AffordabilityTransactionList
		if err = json.Unmarshal(data, &list); err != nil {
			return nil, err
		}
		if len(list.Data) > 0 {
			transactions = append(transactions, list.Data...)
		}
		callURL = list.Links.Next
	}
	return transactions, nil
}
