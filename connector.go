package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type ConnectorList struct {
	Type       string      `json:"type"`
	TotalCount int         `json:"totalCount"`
	Data       []Connector `json:"data"`
	Links      struct {
		Self string `json:"self"`
	} `json:"links"`
}

type Connector struct {
	Type          string `json:"type"`
	ID            string `json:"id"`
	Status        string `json:"status"`
	Method        string `json:"method"`
	Authorization struct {
		Meta struct {
			ForgottenPasswordUrl    string `json:"forgotten_password_url"`
			LoginIdCaption          string `json:"login_id_caption"`
			PasswordCaption         string `json:"password_caption"`
			SecondaryLoginIdCaption string `json:"secondary_login_id_caption"`
			SecurityCodeCaption     string `json:"security_code_caption"`
		} `json:"meta"`
		Type string `json:"type"`
	} `json:"authorization"`
	Institution struct {
		Type      string `json:"type"`
		Name      string `json:"name"`
		Country   string `json:"country"`
		ShortName string `json:"shortName"`
		Tier      string `json:"tier"`
		Logo      struct {
			Colors struct {
				Primary string `json:"primary"`
			} `json:"colors"`
			Links struct {
				Full   string `json:"full"`
				Square string `json:"square"`
			} `json:"links"`
			Type string `json:"type"`
		} `json:"logo"`
	} `json:"institution"`

	Scopes []string `json:"scopes"`
	Stage  string   `json:"stage"`
	Stats  struct {
		AverageDurationMs struct {
			RetrieveAccounts     int `json:"retrieveAccounts"`
			RetrieveMeta         int `json:"retrieveMeta"`
			RetrieveTransactions int `json:"retrieveTransactions"`
			Total                int `json:"total"`
			VerifyCredentials    int `json:"verifyCredentials"`
		} `json:"averageDurationMs"`
	} `json:"stats"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) Connector(ctx context.Context, connectorID, method string) (Connector, error) {
	connector, err := c.connector(ctx, connectorID, method)
	if err != nil && !IsUnauthorizedErr(err) {
		return connector, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return Connector{}, err
	}
	return c.connector(ctx, connectorID, method)
}

func (c *Client) Connectors(ctx context.Context) (ConnectorList, error) {
	connectorList, err := c.connectors(ctx)
	if err != nil && !IsUnauthorizedErr(err) {
		return connectorList, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return ConnectorList{}, err
	}
	return c.connectors(ctx)
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) connector(ctx context.Context, connectorID, method string) (Connector, error) {
	callURL, err := url.JoinPath(baseURL, "connectors", connectorID, method)
	if err != nil {
		return Connector{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return Connector{}, err
	}

	var connector Connector
	return connector, json.Unmarshal(data, &connector)
}

func (c *Client) connectors(ctx context.Context) (ConnectorList, error) {
	callURL, err := url.JoinPath(baseURL, "connectors")
	if err != nil {
		return ConnectorList{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return ConnectorList{}, err
	}

	var list ConnectorList
	return list, json.Unmarshal(data, &list)
}
