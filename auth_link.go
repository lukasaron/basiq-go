package basiq

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type AuthLinkParams struct {
	Mobile string `json:"mobile,omitempty"`
}

type AuthLink struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Mobile    int64  `json:"mobile"`
	UserID    string `json:"userId"`
	ExpiresAt string `json:"expiresAt"`
	Links     struct {
		Public string `json:"public"`
		Self   string `json:"self"`
	} `json:"links"`
}

//---------------------------------------------------------------------------------------------------------------------

func (a *API) AuthLink(ctx context.Context, userID string) (AuthLink, error) {
	authLink, err := a.authLink(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return authLink, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return AuthLink{}, err
	}
	return a.authLink(ctx, userID)
}

func (a *API) CreateAuthLink(ctx context.Context, userID string, params AuthLinkParams) (AuthLink, error) {
	authLink, err := a.createAuthLink(ctx, userID, params)
	if err != nil && !IsUnauthorizedErr(err) {
		return authLink, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return AuthLink{}, err
	}
	return a.createAuthLink(ctx, userID, params)
}

func (a *API) DeleteAuthLink(ctx context.Context, userID string) error {
	err := a.deleteAuthLink(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return err
	}
	if err = a.Authenticate(ctx); err != nil {
		return err
	}
	return a.deleteAuthLink(ctx, userID)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) authLink(ctx context.Context, userID string) (AuthLink, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "auth_link")
	if err != nil {
		return AuthLink{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return AuthLink{}, err
	}

	var link AuthLink
	return link, json.Unmarshal(data, &link)
}

func (a *API) createAuthLink(ctx context.Context, userID string, params AuthLinkParams) (AuthLink, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "auth_link")
	if err != nil {
		return AuthLink{}, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return AuthLink{}, err
	}

	data, err := a.makeCall(ctx, http.MethodPost, callURL, bytes.NewReader(payload))
	if err != nil {
		return AuthLink{}, err
	}

	var link AuthLink
	return link, json.Unmarshal(data, &link)
}

func (a *API) deleteAuthLink(ctx context.Context, userID string) error {
	callURL, err := url.JoinPath(baseURL, "users", userID, "auth_link")
	if err != nil {
		return err
	}

	_, err = a.makeCall(ctx, http.MethodDelete, callURL, nil)
	return err
}
