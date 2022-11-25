package basiq

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type AffordabilityParams struct {
	Accounts  []string `json:"accounts,omitempty"`
	FromMonth string   `json:"fromMonth,omitempty"`
	ToMonth   string   `json:"toMonth,omitempty"`
}

type Affordability struct {
	Type          string `json:"type"`
	ID            string `json:"id"`
	CoverageDays  int    `json:"coverageDays"`
	FromMonth     string `json:"fromMonth"`
	ToMonth       string `json:"toMonth"`
	GeneratedDate string `json:"generatedDate"`
	Assets        []struct {
		Type    string `json:"type"`
		Account struct {
			Product string `json:"product"`
			Type    string `json:"type"`
		} `json:"account"`
		AvailableFunds  string `json:"availableFunds"`
		Balance         string `json:"balance"`
		Currency        string `json:"currency"`
		Institution     string `json:"institution"`
		Previous6Months struct {
			MaxBalance string `json:"maxBalance"`
			MinBalance string `json:"minBalance"`
		} `json:"previous6Months"`
	} `json:"assets"`
	External []struct {
		ChangeHistory []struct {
			Amount string `json:"amount"`
			Date   string `json:"date"`
			Source string `json:"source"`
		} `json:"changeHistory"`
		Payments []struct {
			AmountAvg        string `json:"amountAvg"`
			AmountAvgMonthly string `json:"amountAvgMonthly"`
			First            string `json:"first"`
			Last             string `json:"last"`
			NoOccurrences    int    `json:"noOccurrences"`
			Total            string `json:"total"`
		} `json:"payments"`
		Source string `json:"source"`
	} `json:"external"`
	Liabilities struct {
		Credit []struct {
			Account struct {
				Product string `json:"product"`
				Type    string `json:"type"`
			} `json:"account"`
			AvailableFunds  string `json:"availableFunds"`
			Balance         string `json:"balance"`
			CreditLimit     string `json:"creditLimit"`
			Currency        string `json:"currency"`
			Institution     string `json:"institution"`
			Previous6Months struct {
				CashAdvances float64 `json:"cashAdvances"`
			} `json:"previous6Months"`
			PreviousMonth struct {
				MaxBalance   float64 `json:"maxBalance"`
				MinBalance   float64 `json:"minBalance"`
				TotalCredits string  `json:"totalCredits"`
				TotalDebits  float64 `json:"totalDebits"`
			} `json:"previousMonth"`
		} `json:"credit"`
		Loan []struct {
			Account struct {
				Product string `json:"product"`
				Type    string `json:"type"`
			} `json:"account"`
			AvailableFunds string `json:"availableFunds"`
			Balance        string `json:"balance"`
			ChangeHistory  []struct {
				Amount    string `json:"amount"`
				Date      string `json:"date"`
				Direction string `json:"direction"`
				Source    string `json:"source"`
			} `json:"changeHistory"`
			Currency        string `json:"currency"`
			Institution     string `json:"institution"`
			Previous6Months struct {
				Arrears string `json:"arrears"`
			} `json:"previous6Months"`
			PreviousMonth struct {
				TotalCredits         string  `json:"totalCredits"`
				TotalDebits          float64 `json:"totalDebits"`
				TotalInterestCharged float64 `json:"totalInterestCharged"`
				TotalRepayments      string  `json:"totalRepayments"`
			} `json:"previousMonth"`
		} `json:"loan"`
	} `json:"liabilities"`
	Summary struct {
		Assets                      string  `json:"assets"`
		CreditLimit                 string  `json:"creditLimit"`
		Expenses                    string  `json:"expenses"`
		Liabilities                 string  `json:"liabilities"`
		LoanRepaymentMonthly        string  `json:"loanRepaymentMonthly"`
		NetPosition                 float64 `json:"netPosition"`
		PotentialLiabilitiesMonthly int     `json:"potentialLiabilitiesMonthly"`
		RegularIncome               struct {
			Previous3Months struct {
				AvgMonthly string `json:"avgMonthly"`
			} `json:"previous3Months"`
		} `json:"regularIncome"`
		Savings string `json:"savings"`
	} `json:"summary"`
	Links struct {
		Accounts []string `json:"accounts"`
		Expenses string   `json:"expenses"`
		Income   string   `json:"income"`
		Self     string   `json:"self"`
	} `json:"links"`
}

//---------------------------------------------------------------------------------------------------------------------

func (a *API) Affordability(ctx context.Context, userID, snapshotID string) (Affordability, error) {
	affordability, err := a.affordability(ctx, userID, snapshotID)
	if err != nil && !IsUnauthorizedErr(err) {
		return affordability, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return Affordability{}, err
	}
	return a.affordability(ctx, userID, snapshotID)
}

func (a *API) CreateAffordability(ctx context.Context, userID string, params AffordabilityParams) (Affordability, error) {
	affordability, err := a.createAffordability(ctx, userID, params)
	if err != nil && !IsUnauthorizedErr(err) {
		return affordability, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return Affordability{}, err
	}
	return a.createAffordability(ctx, userID, params)
}

//---------------------------------------------------------------------------------------------------------------------

func (a *API) affordability(ctx context.Context, userID, snapshotID string) (Affordability, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "affordability", snapshotID)
	if err != nil {
		return Affordability{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return Affordability{}, err
	}

	var affordability Affordability
	return affordability, json.Unmarshal(data, &affordability)
}

func (a *API) createAffordability(ctx context.Context, userID string, params AffordabilityParams) (Affordability, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "affordability")
	if err != nil {
		return Affordability{}, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return Affordability{}, err
	}

	data, err := a.makeCall(ctx, http.MethodPost, callURL, bytes.NewReader(payload))
	if err != nil {
		return Affordability{}, err
	}

	var affordability Affordability
	return affordability, json.Unmarshal(data, &affordability)
}
