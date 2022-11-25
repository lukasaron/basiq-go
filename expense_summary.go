package basiq

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type ExpenseSummaryParams struct {
	Accounts  []string `json:"accounts,omitempty"`
	FromMonth string   `json:"fromMonth,omitempty"`
	ToMonth   string   `json:"toMonth,omitempty"`
}

type ExpenseSummary struct {
	Type         string `json:"type"`
	ID           string `json:"id"`
	CoverageDays int    `json:"coverageDays"`
	BankFees     struct {
		AvgMonthly    string `json:"avgMonthly"`
		ChangeHistory []struct {
			Amount string `json:"amount"`
			Date   string `json:"date"`
		} `json:"changeHistory"`
		Summary string `json:"summary"`
	} `json:"bankFees"`
	CashWithdrawals struct {
		AvgMonthly    string `json:"avgMonthly"`
		ChangeHistory []struct {
			Amount string `json:"amount"`
			Date   string `json:"date"`
		} `json:"changeHistory"`
		Summary string `json:"summary"`
	} `json:"cashWithdrawals"`
	ExternalTransfers struct {
		AvgMonthly    string `json:"avgMonthly"`
		ChangeHistory []struct {
			Amount string `json:"amount"`
			Date   string `json:"date"`
		} `json:"changeHistory"`
		Summary string `json:"summary"`
	} `json:"externalTransfers"`
	FromMonth     string `json:"fromMonth"`
	LoanInterests struct {
		AvgMonthly    string `json:"avgMonthly"`
		ChangeHistory []struct {
			Amount string `json:"amount"`
			Date   string `json:"date"`
		} `json:"changeHistory"`
		Summary string `json:"summary"`
	} `json:"loanInterests"`
	LoanRepayments struct {
		AvgMonthly    string `json:"avgMonthly"`
		ChangeHistory []struct {
			Amount string `json:"amount"`
			Date   string `json:"date"`
		} `json:"changeHistory"`
		Summary string `json:"summary"`
	} `json:"loanRepayments"`
	Payments []struct {
		AvgMonthly      string `json:"avgMonthly"`
		Division        string `json:"division"`
		PercentageTotal string `json:"percentageTotal"`
		SubCategory     []struct {
			Category struct {
				ExpenseClass struct {
					ClassCode     string `json:"classCode"`
					ClassTitle    string `json:"classTitle"`
					DivisionCode  string `json:"divisionCode"`
					DivisionTitle string `json:"divisionTitle"`
				} `json:"expenseClass"`
			} `json:"category"`
			ChangeHistory []struct {
				Amount string `json:"amount"`
				Date   string `json:"date"`
			} `json:"changeHistory"`
			Summary string `json:"summary"`
		} `json:"subCategory"`
	} `json:"payments"`
	ToMonth string `json:"toMonth"`
	Links   struct {
		Accounts []string `json:"accounts"`
		Self     string   `json:"self"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) ExpenseSummary(ctx context.Context, userID, snapshotID string) (ExpenseSummary, error) {
	expenseSummary, err := c.expenseSummary(ctx, userID, snapshotID)
	if err != nil && !IsUnauthorizedErr(err) {
		return expenseSummary, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return ExpenseSummary{}, err
	}
	return c.expenseSummary(ctx, userID, snapshotID)
}

func (c *Client) CreateExpenseSummary(ctx context.Context, userID string, params ExpenseSummaryParams) (ExpenseSummary, error) {
	expenseSummary, err := c.createExpenseSummary(ctx, userID, params)
	if err != nil && !IsUnauthorizedErr(err) {
		return expenseSummary, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return ExpenseSummary{}, err
	}
	return c.createExpenseSummary(ctx, userID, params)
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) expenseSummary(ctx context.Context, userID, snapshotID string) (ExpenseSummary, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "expenses", snapshotID)
	if err != nil {
		return ExpenseSummary{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return ExpenseSummary{}, err
	}

	var summary ExpenseSummary
	return summary, json.Unmarshal(data, &summary)
}

func (c *Client) createExpenseSummary(ctx context.Context, userID string, params ExpenseSummaryParams) (ExpenseSummary, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "expenses")
	if err != nil {
		return ExpenseSummary{}, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return ExpenseSummary{}, err
	}

	data, err := c.makeCall(ctx, http.MethodPost, callURL, bytes.NewReader(payload))
	if err != nil {
		return ExpenseSummary{}, err
	}

	var summary ExpenseSummary
	return summary, json.Unmarshal(data, &summary)
}
