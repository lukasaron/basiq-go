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

func (c *Client) AuthLink(ctx context.Context, userID string) (AuthLink, error) {
	authLink, err := c.authLink(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return authLink, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return AuthLink{}, err
	}
	return c.authLink(ctx, userID)
}

func (c *Client) CreateAuthLink(ctx context.Context, userID string, params AuthLinkParams) (AuthLink, error) {
	authLink, err := c.createAuthLink(ctx, userID, params)
	if err != nil && !IsUnauthorizedErr(err) {
		return authLink, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return AuthLink{}, err
	}
	return c.createAuthLink(ctx, userID, params)
}

func (c *Client) DeleteAuthLink(ctx context.Context, userID string) error {
	err := c.deleteAuthLink(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return err
	}
	if err = c.Authenticate(ctx); err != nil {
		return err
	}
	return c.deleteAuthLink(ctx, userID)
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) authLink(ctx context.Context, userID string) (AuthLink, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "auth_link")
	if err != nil {
		return AuthLink{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return AuthLink{}, err
	}

	var link AuthLink
	return link, json.Unmarshal(data, &link)
}

func (c *Client) createAuthLink(ctx context.Context, userID string, params AuthLinkParams) (AuthLink, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "auth_link")
	if err != nil {
		return AuthLink{}, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return AuthLink{}, err
	}

	data, err := c.makeCall(ctx, http.MethodPost, callURL, bytes.NewReader(payload))
	if err != nil {
		return AuthLink{}, err
	}

	var link AuthLink
	return link, json.Unmarshal(data, &link)
}

func (c *Client) deleteAuthLink(ctx context.Context, userID string) error {
	callURL, err := url.JoinPath(baseURL, "users", userID, "auth_link")
	if err != nil {
		return err
	}

	_, err = c.makeCall(ctx, http.MethodDelete, callURL, nil)
	return err
}
