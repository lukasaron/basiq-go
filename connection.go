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

func (a *API) Connection(ctx context.Context, userID, connectionID string) (Connection, error) {
	connection, err := a.connection(ctx, userID, connectionID)
	if err != nil && !IsUnauthorizedErr(err) {
		return connection, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return Connection{}, err
	}
	return a.connection(ctx, userID, connectionID)
}

func (a *API) Connections(ctx context.Context, userID string) (ConnectionList, error) {
	connectionList, err := a.connections(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return connectionList, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return ConnectionList{}, err
	}
	return a.connections(ctx, userID)
}

func (a *API) RefreshConnection(ctx context.Context, userID, connectionID string) (Connection, error) {
	connection, err := a.refreshConnection(ctx, userID, connectionID)
	if err != nil && !IsUnauthorizedErr(err) {
		return connection, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return Connection{}, err
	}
	return a.refreshConnection(ctx, userID, connectionID)
}

func (a *API) RefreshConnections(ctx context.Context, userID string) (ConnectionList, error) {
	connectionList, err := a.refreshConnections(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return connectionList, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return ConnectionList{}, err
	}
	return a.refreshConnections(ctx, userID)
}

func (a *API) DeleteConnection(ctx context.Context, userID, connectionID string) error {
	err := a.deleteConnection(ctx, userID, connectionID)
	if err != nil && !IsUnauthorizedErr(err) {
		return err
	}
	if err = a.Authenticate(ctx); err != nil {
		return err
	}
	return a.deleteConnection(ctx, userID, connectionID)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) connection(ctx context.Context, userID, connectionID string) (Connection, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "connections", connectionID)
	if err != nil {
		return Connection{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return Connection{}, err
	}

	var connection Connection
	return connection, json.Unmarshal(data, &connection)
}

func (a *API) connections(ctx context.Context, userID string) (ConnectionList, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "connections")
	if err != nil {
		return ConnectionList{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return ConnectionList{}, err
	}

	var list ConnectionList
	return list, json.Unmarshal(data, &list)
}

func (a *API) refreshConnection(ctx context.Context, userID, connectionID string) (Connection, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "connections", connectionID, "refresh")
	if err != nil {
		return Connection{}, err
	}

	data, err := a.makeCall(ctx, http.MethodPost, callURL, nil)
	if err != nil {
		return Connection{}, err
	}

	var connection Connection
	return connection, json.Unmarshal(data, &connection)
}

func (a *API) refreshConnections(ctx context.Context, userID string) (ConnectionList, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "connections", "refresh")
	if err != nil {
		return ConnectionList{}, err
	}

	data, err := a.makeCall(ctx, http.MethodPost, callURL, nil)
	if err != nil {
		return ConnectionList{}, err
	}

	var list ConnectionList
	return list, json.Unmarshal(data, &list)
}

func (a *API) deleteConnection(ctx context.Context, userID, connectionID string) error {
	callURL, err := url.JoinPath(baseURL, "users", userID, "connections", connectionID)
	if err != nil {
		return err
	}

	_, err = a.makeCall(ctx, http.MethodDelete, callURL, nil)
	return err
}
