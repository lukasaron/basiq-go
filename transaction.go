package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type TransactionList struct {
	Type  string        `json:"type"`
	Count int           `json:"count"`
	Size  int           `json:"size"`
	Data  []Transaction `json:"data"`
	Links struct {
		Self string `json:"self"`
		Next string `json:"next"`
	} `json:"links"`
}

type Transaction struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Account     string `json:"account"`
	Amount      string `json:"amount"`
	Balance     string `json:"balance"`
	Class       string `json:"class"`
	Connection  string `json:"connection"`
	Description string `json:"description"`
	Direction   string `json:"direction"`
	Enrich      struct {
		Category struct {
			Anzsic struct {
				Class struct {
					Title string `json:"title"`
					Code  string `json:"code"`
				} `json:"class"`
				Division struct {
					Title string `json:"title"`
					Code  string `json:"code"`
				} `json:"division"`
				Group struct {
					Code  string `json:"code"`
					Title string `json:"title"`
				} `json:"group"`
				Subdivision struct {
					Code  string `json:"code"`
					Title string `json:"title"`
				} `json:"subdivision"`
			} `json:"anzsic"`
		} `json:"category"`
		Location struct {
			Country          string `json:"country"`
			FormattedAddress string `json:"formattedAddress"`
			Geometry         struct {
				Lat string `json:"lat"`
				Lng string `json:"lng"`
			} `json:"geometry"`
			PostalCode string `json:"postalCode"`
			Route      string `json:"route"`
			RouteNo    string `json:"routeNo"`
			State      string `json:"state"`
			Suburb     string `json:"suburb"`
		} `json:"location"`
		Merchant struct {
			ID           string `json:"id"`
			BusinessName string `json:"businessName"`
			ABN          int64  `json:"ABN"`
			LogoMaster   string `json:"logoMaster"`
			LogoThumb    string `json:"logoThumb"`
			PhoneNumber  struct {
				International string `json:"international"`
				Local         string `json:"local"`
			} `json:"phoneNumber"`
			Website string `json:"website"`
		} `json:"merchant"`
	} `json:"enrich"`
	Institution     string `json:"institution"`
	PostDate        string `json:"postDate"`
	Status          string `json:"status"`
	TransactionDate string `json:"transactionDate"`
	Links           struct {
		Account     string `json:"account"`
		Institution string `json:"institution"`
		Self        string `json:"self"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) Transaction(ctx context.Context, userID, transactionID string) (Transaction, error) {
	transaction, err := a.transaction(ctx, userID, transactionID)
	if err != nil && !IsUnauthorizedErr(err) {
		return transaction, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return Transaction{}, err
	}
	return a.transaction(ctx, userID, transactionID)
}

func (a *API) Transactions(ctx context.Context, userID string) ([]Transaction, error) {
	transactions, err := a.transactions(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return transactions, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return nil, err
	}
	return a.transactions(ctx, userID)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) transaction(ctx context.Context, userID, transactionID string) (Transaction, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "transactions", transactionID)
	if err != nil {
		return Transaction{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return Transaction{}, err
	}

	var transaction Transaction
	return transaction, json.Unmarshal(data, &transaction)
}

func (a *API) transactions(ctx context.Context, userID string) ([]Transaction, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "transactions")
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	for callURL != "" {
		data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
		if err != nil {
			return nil, err
		}
		var list TransactionList
		if err = json.Unmarshal(data, &list); err != nil {
			return nil, err
		}
		if len(list.Data) > 0 {
			transactions = append(transactions, list.Data...)
		}
		callURL = list.Links.Next
	}

	return transactions, err
}
