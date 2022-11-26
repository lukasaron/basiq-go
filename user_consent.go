package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

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

func (a *API) UserConsent(ctx context.Context, userID string) (UserConsent, error) {
	userConsent, err := a.userConsent(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return userConsent, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return UserConsent{}, err
	}
	return a.userConsent(ctx, userID)
}

func (a *API) DeleteUserConsent(ctx context.Context, userID, consentID string) error {
	err := a.deleteUserConsent(ctx, userID, consentID)
	if err != nil && !IsUnauthorizedErr(err) {
		return err
	}
	if err = a.Authenticate(ctx); err != nil {
		return err
	}
	return a.deleteUserConsent(ctx, userID, consentID)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) userConsent(ctx context.Context, userID string) (UserConsent, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "consents")
	if err != nil {
		return UserConsent{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return UserConsent{}, err
	}

	var consent UserConsent
	return consent, json.Unmarshal(data, &consent)
}

func (a *API) deleteUserConsent(ctx context.Context, userID, consentID string) error {
	callURL, err := url.JoinPath(baseURL, "users", userID, "consents", consentID)
	if err != nil {
		return err
	}

	_, err = a.makeCall(ctx, http.MethodDelete, callURL, nil)
	return err
}
