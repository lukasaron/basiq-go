package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type ConnectionList struct {
	Type  string       `json:"type"`
	Data  []Connection `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type Connection struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Method      string `json:"method"`
	CreatedDate string `json:"createdDate"`
	LastUsed    string `json:"lastUsed"`
	Status      string `json:"status"`
	Accounts    struct {
		Type string `json:"type"`
		Data []struct {
			Type     string `json:"type"`
			ID       string `json:"id"`
			Name     string `json:"name"`
			Currency string `json:"currency"`
			Class    struct {
				Type    string `json:"type"`
				Product string `json:"product"`
			} `json:"class"`
			AccountNo      string `json:"accountNo"`
			AvailableFunds string `json:"availableFunds"`
			Balance        string `json:"balance"`
			LastUpdated    string `json:"lastUpdated"`
			Status         string `json:"status"`
			Links          struct {
				Transactions string `json:"transactions"`
				Self         string `json:"self"`
			} `json:"links"`
		} `json:"data"`
	} `json:"accounts"`
	Institution struct {
		Id    string `json:"id"`
		Type  string `json:"type"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"institution"`
	Profile struct {
		EmailAddresses    []string `json:"emailAddresses"`
		FirstName         string   `json:"firstName"`
		FullName          string   `json:"fullName"`
		LastName          string   `json:"lastName"`
		MiddleName        string   `json:"middleName"`
		PhoneNumbers      []string `json:"phoneNumbers"`
		PhysicalAddresses []struct {
			AddressLine1     string `json:"addressLine1"`
			AddressLine2     string `json:"addressLine2"`
			AddressLine3     string `json:"addressLine3"`
			City             string `json:"city"`
			Country          string `json:"country"`
			CountryCode      string `json:"countryCode"`
			FormattedAddress string `json:"formattedAddress"`
			Postcode         string `json:"postcode"`
			State            string `json:"state"`
		} `json:"physicalAddresses"`
	} `json:"profile"`
	Links struct {
		Accounts     string `json:"accounts"`
		Self         string `json:"self"`
		Transactions string `json:"transactions"`
		User         string `json:"user"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) Connection(ctx context.Context, userID, connectionID string) (Connection, error) {
	connection, err := c.connection(ctx, userID, connectionID)
	if err != nil && !IsUnauthorizedErr(err) {
		return connection, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return Connection{}, err
	}
	return c.connection(ctx, userID, connectionID)
}

func (c *Client) Connections(ctx context.Context, userID string) (ConnectionList, error) {
	connectionList, err := c.connections(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return connectionList, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return ConnectionList{}, err
	}
	return c.connections(ctx, userID)
}

func (c *Client) RefreshConnection(ctx context.Context, userID, connectionID string) (Connection, error) {
	connection, err := c.refreshConnection(ctx, userID, connectionID)
	if err != nil && !IsUnauthorizedErr(err) {
		return connection, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return Connection{}, err
	}
	return c.refreshConnection(ctx, userID, connectionID)
}

func (c *Client) RefreshConnections(ctx context.Context, userID string) (ConnectionList, error) {
	connectionList, err := c.refreshConnections(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return connectionList, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return ConnectionList{}, err
	}
	return c.refreshConnections(ctx, userID)
}

func (c *Client) DeleteConnection(ctx context.Context, userID, connectionID string) error {
	err := c.deleteConnection(ctx, userID, connectionID)
	if err != nil && !IsUnauthorizedErr(err) {
		return err
	}
	if err = c.Authenticate(ctx); err != nil {
		return err
	}
	return c.deleteConnection(ctx, userID, connectionID)
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) connection(ctx context.Context, userID, connectionID string) (Connection, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "connections", connectionID)
	if err != nil {
		return Connection{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return Connection{}, err
	}

	var connection Connection
	return connection, json.Unmarshal(data, &connection)
}

func (c *Client) connections(ctx context.Context, userID string) (ConnectionList, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "connections")
	if err != nil {
		return ConnectionList{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return ConnectionList{}, err
	}

	var list ConnectionList
	return list, json.Unmarshal(data, &list)
}

func (c *Client) refreshConnection(ctx context.Context, userID, connectionID string) (Connection, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "connections", connectionID, "refresh")
	if err != nil {
		return Connection{}, err
	}

	data, err := c.makeCall(ctx, http.MethodPost, callURL, nil)
	if err != nil {
		return Connection{}, err
	}

	var connection Connection
	return connection, json.Unmarshal(data, &connection)
}

func (c *Client) refreshConnections(ctx context.Context, userID string) (ConnectionList, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "connections", "refresh")
	if err != nil {
		return ConnectionList{}, err
	}

	data, err := c.makeCall(ctx, http.MethodPost, callURL, nil)
	if err != nil {
		return ConnectionList{}, err
	}

	var list ConnectionList
	return list, json.Unmarshal(data, &list)
}

func (c *Client) deleteConnection(ctx context.Context, userID, connectionID string) error {
	callURL, err := url.JoinPath(baseURL, "users", userID, "connections", connectionID)
	if err != nil {
		return err
	}

	_, err = c.makeCall(ctx, http.MethodDelete, callURL, nil)
	return err
}
