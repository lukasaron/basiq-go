package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type UserConsentList struct {
	Type string        `json:"type"`
	Size int           `json:"size"`
	Data []UserConsent `json:"data"`
}

type UserConsent struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Created    string `json:"created"`
	Updated    string `json:"updated"`
	ExpiryDate string `json:"expiryDate"`
	Status     string `json:"status"`
	Purpose    struct {
		Primary struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"primary"`
	} `json:"purpose"`
	Data struct {
		RetainData  bool `json:"retainData"`
		Permissions []struct {
			Scope       string `json:"scope"`
			Required    bool   `json:"required"`
			Entity      string `json:"entity"`
			Information struct {
				Name          string   `json:"name"`
				Description   string   `json:"description"`
				AttributeList []string `json:"attributeList"`
			} `json:"information"`
			Purpose struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"purpose"`
		} `json:"permissions"`
	} `json:"data"`
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) UserConsent(ctx context.Context, userID string) (UserConsent, error) {
	userConsent, err := c.userConsent(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return userConsent, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return UserConsent{}, err
	}
	return c.userConsent(ctx, userID)
}

func (c *Client) DeleteUserConsent(ctx context.Context, userID, consentID string) error {
	err := c.deleteUserConsent(ctx, userID, consentID)
	if err != nil && !IsUnauthorizedErr(err) {
		return err
	}
	if err = c.Authenticate(ctx); err != nil {
		return err
	}
	return c.deleteUserConsent(ctx, userID, consentID)
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) userConsent(ctx context.Context, userID string) (UserConsent, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "consents")
	if err != nil {
		return UserConsent{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return UserConsent{}, err
	}

	var consent UserConsent
	return consent, json.Unmarshal(data, &consent)
}

func (c *Client) deleteUserConsent(ctx context.Context, userID, consentID string) error {
	callURL, err := url.JoinPath(baseURL, "users", userID, "consents", consentID)
	if err != nil {
		return err
	}

	_, err = c.makeCall(ctx, http.MethodDelete, callURL, nil)
	return err
}
