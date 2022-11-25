package basiq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AuthScope string

var (
	ClientScope AuthScope = "CLIENT_ACCESS"
	ServerScope AuthScope = "SERVER_ACCESS"
)

type AuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) Authenticate(ctx context.Context) error {
	c.m.Lock()
	defer c.m.Unlock()

	// when multiple threads want's to get a token the mutex allows only one enter into the section and
	// this checks whether the last updated happen within a pause timeframe. If yes the token has been already refreshed
	// otherwise refresh the token.
	if c.authorizedAt.Sub(time.Now()).Abs().Seconds() < defaultAuthPauseSec {
		return nil
	}

	authToken, err := c.authToken(ctx, c.apiKey, c.scope, c.userID)
	if err != nil {
		return err
	}

	c.headers.Set("Authorization", fmt.Sprintf("Bearer %s", authToken.AccessToken))
	c.authorizedAt = time.Now()

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (c *Client) authToken(ctx context.Context, apiKey string, scope AuthScope, userID string) (AuthToken, error) {
	callURL, err := url.JoinPath(baseURL, "token")
	if err != nil {
		return AuthToken{}, err
	}

	c.headers.Set("Authorization", fmt.Sprintf("Basic %s", apiKey))

	urlValues := url.Values{
		"scope": {string(scope)},
	}
	if userID != "" {
		urlValues.Set("userId", userID)
	}

	data, err := c.makeCall(ctx, http.MethodPost, callURL, strings.NewReader(urlValues.Encode()))
	if err != nil {
		return AuthToken{}, err
	}

	var token AuthToken
	return token, json.Unmarshal(data, &token)
}
