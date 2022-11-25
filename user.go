package basiq

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type User struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Name      string `json:"name"`
	Accounts  struct {
		Type  string `json:"type"`
		Count int    `json:"count"`
		Data  []struct {
			Type  string `json:"type"`
			ID    string `json:"id"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
	} `json:"accounts"`
	Connections struct {
		Type  string `json:"type"`
		Count int    `json:"count"`
		Data  []struct {
			Type  string `json:"type"`
			ID    string `json:"id"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
	} `json:"connections"`
	Links struct {
		Accounts     string `json:"accounts"`
		Connections  string `json:"connections"`
		Self         string `json:"self"`
		Transactions string `json:"transactions"`
	} `json:"links"`
}

type UserParams struct {
	Email     string `json:"email,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

//---------------------------------------------------------------------------------------------------------------------

func (c *Client) User(ctx context.Context, userID string) (User, error) {
	user, err := c.user(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return user, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return User{}, err
	}
	return c.user(ctx, userID)
}

func (c *Client) CreateUser(ctx context.Context, params UserParams) (User, error) {
	user, err := c.createUser(ctx, params)
	if err != nil && !IsUnauthorizedErr(err) {
		return user, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return User{}, err
	}
	return c.createUser(ctx, params)
}

func (c *Client) UpdateUser(ctx context.Context, userID string, params UserParams) (User, error) {
	user, err := c.updateUser(ctx, userID, params)
	if err != nil && !IsUnauthorizedErr(err) {
		return user, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return User{}, err
	}
	return c.updateUser(ctx, userID, params)
}

func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	err := c.deleteUser(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return err
	}
	if err = c.Authenticate(ctx); err != nil {
		return err
	}
	return c.deleteUser(ctx, userID)
}

//---------------------------------------------------------------------------------------------------------------------

func (c *Client) user(ctx context.Context, userID string) (User, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID)
	if err != nil {
		return User{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return User{}, err
	}

	var user User
	return user, json.Unmarshal(data, &user)
}

func (c *Client) createUser(ctx context.Context, params UserParams) (User, error) {
	callURL, err := url.JoinPath(baseURL, "users")
	if err != nil {
		return User{}, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return User{}, err
	}

	data, err := c.makeCall(ctx, http.MethodPost, callURL, bytes.NewReader(payload))
	if err != nil {
		return User{}, err
	}

	var user User
	return user, json.Unmarshal(data, &user)
}

func (c *Client) updateUser(ctx context.Context, userID string, params UserParams) (User, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID)
	if err != nil {
		return User{}, err
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return User{}, err
	}

	data, err := c.makeCall(ctx, http.MethodPost, callURL, bytes.NewReader(payload))
	if err != nil {
		return User{}, err
	}

	var user User
	return user, json.Unmarshal(data, &user)
}

func (c *Client) deleteUser(ctx context.Context, userID string) error {
	callURL, err := url.JoinPath(baseURL, "users", userID)
	if err != nil {
		return err
	}

	_, err = c.makeCall(ctx, http.MethodDelete, callURL, nil)
	return err
}
