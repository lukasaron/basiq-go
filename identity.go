package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type IdentityList struct {
	Type  string     `json:"type"`
	Count int        `json:"count"`
	Data  []Identity `json:"data"`
}

type Identity struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Links   struct {
		Self string `json:"self"`
		Job  string `json:"job"`
	} `json:"links"`
	Source                string   `json:"source"`
	FullName              string   `json:"fullName"`
	FirstName             string   `json:"firstName"`
	LastName              string   `json:"lastName"`
	MiddleName            string   `json:"middleName"`
	Title                 string   `json:"title"`
	DOB                   string   `json:"DOB"`
	OccupationCode        string   `json:"occupationCode"`
	OccupationCodeVersion string   `json:"occupationCodeVersion"`
	PhoneNumbers          []string `json:"phoneNumbers"`
	Emails                []string `json:"emails"`
	PhysicalAddresses     []struct {
		Type             string `json:"type"`
		AddressLine1     string `json:"addressLine1"`
		AddressLine2     string `json:"addressLine2"`
		AddressLine3     string `json:"addressLine3"`
		Postcode         string `json:"postcode"`
		City             string `json:"city"`
		State            string `json:"state"`
		Country          string `json:"country"`
		CountryCode      string `json:"countryCode"`
		FormattedAddress string `json:"formattedAddress"`
	} `json:"physicalAddresses"`
	Organisation struct {
		AgentFirstName      string `json:"agentFirstName"`
		AgentLastName       string `json:"agentLastName"`
		AgentRole           string `json:"agentRole"`
		BusinessName        string `json:"businessName"`
		LegalName           string `json:"legalName"`
		ShortName           string `json:"shortName"`
		ABN                 string `json:"abn"`
		ACN                 string `json:"acn"`
		IsACNCRegistered    bool   `json:"isACNCRegistered"`
		IndustryCode        string `json:"industryCode"`
		IndustryCodeVersion string `json:"industryCodeVersion"`
		OrganisationType    string `json:"organisationType"`
		RegisteredCountry   string `json:"registeredCountry"`
	} `json:"organisation"`
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) Identity(ctx context.Context, userID, identityID string) (Identity, error) {
	identity, err := a.identity(ctx, userID, identityID)
	if err != nil && !IsUnauthorizedErr(err) {
		return identity, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return Identity{}, err
	}
	return a.identity(ctx, userID, identityID)
}

func (a *API) Identities(ctx context.Context, userID string) ([]Identity, error) {
	identities, err := a.identities(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return identities, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return nil, err
	}
	return a.identities(ctx, userID)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) identity(ctx context.Context, userID, identityID string) (Identity, error) {
	callURl, err := url.JoinPath(baseURL, "users", userID, "identities", identityID)
	if err != nil {
		return Identity{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURl, nil)
	if err != nil {
		return Identity{}, err
	}

	var identity Identity
	return identity, json.Unmarshal(data, &identity)
}

func (a *API) identities(ctx context.Context, userID string) ([]Identity, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "identities")
	if err != nil {
		return nil, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return nil, err
	}

	var list IdentityList
	if err = json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list.Data, nil
}
