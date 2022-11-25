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

func (a *API) Authenticate(ctx context.Context) error {
	a.m.Lock()
	defer a.m.Unlock()

	// when multiple threads want's to get a token the mutex allows only one enter into the section and
	// this checks whether the last updated happen within a pause timeframe. If yes the token has been already refreshed
	// otherwise refresh the token.
	if a.authorizedAt.Sub(time.Now()).Abs().Seconds() < defaultAuthPauseSec {
		return nil
	}

	authToken, err := a.authToken(ctx, a.apiKey, a.scope, a.userID)
	if err != nil {
		return err
	}

	a.headers.Set("Authorization", fmt.Sprintf("Bearer %s", authToken.AccessToken))
	a.authorizedAt = time.Now()

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (a *API) authToken(ctx context.Context, apiKey string, scope AuthScope, userID string) (AuthToken, error) {
	callURL, err := url.JoinPath(baseURL, "token")
	if err != nil {
		return AuthToken{}, err
	}

	a.headers.Set("Authorization", fmt.Sprintf("Basic %s", apiKey))

	urlValues := url.Values{
		"scope": {string(scope)},
	}
	if userID != "" {
		urlValues.Set("userId", userID)
	}

	data, err := a.makeCall(ctx, http.MethodPost, callURL, strings.NewReader(urlValues.Encode()))
	if err != nil {
		return AuthToken{}, err
	}

	var token AuthToken
	return token, json.Unmarshal(data, &token)
}
