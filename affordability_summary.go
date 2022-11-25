package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type AffordabilitySummaryList struct {
	Type  string                 `json:"type"`
	Data  []AffordabilitySummary `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type AffordabilitySummary struct {
	Type          string `json:"type"`
	ID            string `json:"id"`
	CoverageDays  int    `json:"coverageDays"`
	FromMonth     string `json:"fromMonth"`
	ToMonth       string `json:"toMonth"`
	GeneratedDate string `json:"generatedDate"`
	Institutions  string `json:"institutions"`
	Links         struct {
		Expenses string `json:"expenses"`
		Income   string `json:"income"`
		Self     string `json:"self"`
	} `json:"links"`
}

//---------------------------------------------------------------------------------------------------------------------

func (a *API) AffordabilitySummaries(ctx context.Context, userID string) (AffordabilitySummaryList, error) {
	affordabilitySummaryList, err := a.affordabilitySummaries(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return affordabilitySummaryList, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return AffordabilitySummaryList{}, err
	}
	return a.affordabilitySummaries(ctx, userID)
}

//---------------------------------------------------------------------------------------------------------------------

func (a *API) affordabilitySummaries(ctx context.Context, userID string) (AffordabilitySummaryList, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "affordability")
	if err != nil {
		return AffordabilitySummaryList{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return AffordabilitySummaryList{}, err
	}

	var list AffordabilitySummaryList
	return list, json.Unmarshal(data, &list)
}
