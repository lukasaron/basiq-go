package basiq

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type MFAParams struct {
	MFAResponse []string `json:"mfa-response"`
}

type MFA struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) CreateMFAResponse(ctx context.Context, jobID string, params MFAParams) (MFA, error) {
	mfa, err := a.createMFA(ctx, jobID, params)
	if err != nil && !IsUnauthorizedErr(err) {
		return mfa, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return MFA{}, err
	}
	return a.createMFA(ctx, jobID, params)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) createMFA(ctx context.Context, jobID string, params MFAParams) (MFA, error) {
	callURL, err := url.JoinPath(baseURL, "jobs", jobID, "mfa")
	if err != nil {
		return MFA{}, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return MFA{}, err
	}

	data, err := a.makeCall(ctx, http.MethodPost, callURL, bytes.NewReader(payload))
	if err != nil {
		return MFA{}, err
	}

	var mfa MFA
	return mfa, json.Unmarshal(data, &mfa)
}
