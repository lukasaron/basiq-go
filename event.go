package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type EventList struct {
	Type  string  `json:"type"`
	Data  []Event `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type Event struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	CreatedDate string `json:"createdDate"`
	Entity      string `json:"entity"`
	EventType   string `json:"eventType"`
	UserId      string `json:"userId"`
	DataRef     string `json:"dataRef"`
	Data        []struct {
		Email string `json:"email"`
		ID    string `json:"id"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Mobile string `json:"mobile"`
		Type   string `json:"type"`
	} `json:"data"`
	Links []struct {
		Self string `json:"self"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) Events(ctx context.Context) ([]Event, error) {
	events, err := a.events(ctx)
	if err != nil && !IsUnauthorizedErr(err) {
		return events, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return nil, err
	}
	return a.events(ctx)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) events(ctx context.Context) ([]Event, error) {
	callURL, err := url.JoinPath(baseURL, "events")
	if err != nil {
		return nil, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return nil, err
	}

	var list EventList
	if err = json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list.Data, nil
}
