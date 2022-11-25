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

func (a *API) Events(ctx context.Context) (EventList, error) {
	eventList, err := a.events(ctx)
	if err != nil && !IsUnauthorizedErr(err) {
		return eventList, err
	}
	if err = a.Authenticate(ctx); err != nil {
		return EventList{}, err
	}
	return a.events(ctx)
}

// --------------------------------------------------------------------------------------------------------------------

func (a *API) events(ctx context.Context) (EventList, error) {
	callURL, err := url.JoinPath(baseURL, "events")
	if err != nil {
		return EventList{}, err
	}

	data, err := a.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return EventList{}, err
	}

	var list EventList
	return list, json.Unmarshal(data, &list)
}
