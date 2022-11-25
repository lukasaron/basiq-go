package basiq

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type IncomeSummaryParams struct {
	Accounts  []string `json:"accounts,omitempty"`
	FromMonth string   `json:"fromMonth,omitempty"`
	ToMonth   string   `json:"toMonth,omitempty"`
}

type IncomeSummary struct {
	Type         string `json:"type"`
	ID           string `json:"id"`
	CoverageDays int    `json:"coverageDays"`
	FromMonth    string `json:"fromMonth"`
	ToMonth      string `json:"toMonth"`
	Regular      []struct {
		Source        string `json:"source"`
		AgeDays       int    `json:"ageDays"`
		ChangeHistory []struct {
			Amount string `json:"amount"`
			Date   string `json:"date"`
			Source string `json:"source"`
		} `json:"changeHistory"`
		Current struct {
			Amount   string `json:"amount"`
			Date     string `json:"date"`
			NextDate string `json:"nextDate"`
		} `json:"current"`
		Frequency    string `json:"frequency"`
		Irregularity struct {
			Gaps      []string `json:"gaps"`
			Stability string   `json:"stability"`
		} `json:"irregularity"`
		Previous3Months struct {
			AmountAvg        string `json:"amountAvg"`
			AmountAvgMonthly string `json:"amountAvgMonthly"`
			Variance         string `json:"variance"`
		} `json:"previous3Months"`
	} `json:"regular"`
	Irregular []struct {
		AgeDays              int    `json:"ageDays"`
		AmountAvg            string `json:"amountAvg"`
		AvgMonthlyOccurrence string `json:"avgMonthlyOccurence"`
		ChangeHistory        []struct {
			Amount string `json:"amount"`
			Date   string `json:"date"`
			Source string `json:"source"`
		} `json:"changeHistory"`
		Current struct {
			Amount string `json:"amount"`
			Date   string `json:"date"`
		} `json:"current"`
		Frequency     string `json:"frequency"`
		NoOccurrences int    `json:"noOccurrences"`
		Source        string `json:"source"`
	} `json:"irregular"`
	OtherCredit []struct {
		AgeDays             int    `json:"ageDays"`
		AmountAvg           string `json:"amountAvg"`
		AvgMonthlyOccurence string `json:"avgMonthlyOccurence"`
		ChangeHistory       []struct {
			Amount string `json:"amount"`
			Date   string `json:"date"`
			Source string `json:"source"`
		} `json:"changeHistory"`
		Current struct {
			Amount           string `json:"amount"`
			Date             string `json:"date"`
			OtherCreditLabel string `json:"otherCreditLabel"`
		} `json:"current"`
		Frequency     string `json:"frequency"`
		NoOccurrences int    `json:"noOccurrences"`
		Source        string `json:"source"`
	} `json:"otherCredit"`
	Summary struct {
		IrregularIncomeAvg string      `json:"irregularIncomeAvg"`
		RegularIncomeAvg   string      `json:"regularIncomeAvg"`
		RegularIncomeYTD   string      `json:"regularIncomeYTD"`
		RegularIncomeYear  interface{} `json:"regularIncomeYear"`
	} `json:"summary"`
	Links struct {
		Accounts []string `json:"accounts"`
		Self     string   `json:"self"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) IncomeSummary(ctx context.Context, userID, snapshotID string) (IncomeSummary, error) {
	incomeSummary, err := c.incomeSummary(ctx, userID, snapshotID)
	if err != nil && !IsUnauthorizedErr(err) {
		return incomeSummary, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return IncomeSummary{}, err
	}
	return c.incomeSummary(ctx, userID, snapshotID)
}

func (c *Client) CreateIncomeSummary(ctx context.Context, userID string, params IncomeSummaryParams) (IncomeSummary, error) {
	incomeSummary, err := c.createIncomeSummary(ctx, userID, params)
	if err != nil && !IsUnauthorizedErr(err) {
		return incomeSummary, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return IncomeSummary{}, err
	}
	return c.createIncomeSummary(ctx, userID, params)
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) incomeSummary(ctx context.Context, userID, snapshot string) (IncomeSummary, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "income", snapshot)
	if err != nil {
		return IncomeSummary{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return IncomeSummary{}, err
	}

	var summary IncomeSummary
	return summary, json.Unmarshal(data, &summary)
}

func (c *Client) createIncomeSummary(ctx context.Context, userID string, params IncomeSummaryParams) (IncomeSummary, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "income")
	if err != nil {
		return IncomeSummary{}, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return IncomeSummary{}, err
	}

	data, err := c.makeCall(ctx, http.MethodPost, callURL, bytes.NewReader(payload))
	if err != nil {
		return IncomeSummary{}, err
	}

	var summary IncomeSummary
	return summary, json.Unmarshal(data, &summary)
}
