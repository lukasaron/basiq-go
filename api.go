package basiq

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL = "https://au-api.basiq.io"
)

var defaultAuthPauseSec float64 = 5

var (
	defaultHeaders = http.Header{
		"Accept":        []string{"application/json"},
		"basiq-version": []string{"3.0"},
	}
)

// Config represents the set of required input parameters to start the Client.
type Config struct {
	APIKey string
	Scope  AuthScope
	UserID string
}

// Validate checks all necessary input parameters and returns error when some of them are not set.
func (c Config) Validate() error {
	switch {
	case c.APIKey == "":
		return errors.New("basiq API Key is required")
	case c.Scope == "":
		return errors.New("basic scope is required")
	case c.Scope == ClientScope && c.UserID == "":
		return errors.New("basiq userID is required when CLIENT_ACCESS scope is used")
	default:
		return nil
	}
}

// API is the center logic of the basiq ecosystem. You should crate new instance of the client via calling
// the NewAPI method where all setup and validation of input happens.
// API is thread safe struct.
type API struct {
	apiKey       string
	scope        AuthScope
	userID       string
	authorizedAt time.Time
	headers      http.Header
	m            sync.Mutex
}

// NewAPI instantiates the Client struct and checks all input parameters.
func NewAPI(config Config) (*API, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &API{
		apiKey:  config.APIKey,
		scope:   config.Scope,
		userID:  config.UserID,
		m:       sync.Mutex{},
		headers: defaultHeaders.Clone(),
	}, nil
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) makeCall(ctx context.Context, HTTPMethod, callURL string, payload io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, HTTPMethod, callURL, payload)
	if err != nil {
		return nil, err
	}
	req.Header = a.headers

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return body, nil
	default:
		e := &Error{HttpCode: res.StatusCode}
		_ = json.Unmarshal(body, e)
		return body, e
	}
}
