package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type UserJobList struct {
	Type  string    `json:"type"`
	Size  int       `json:"size"`
	Data  []UserJob `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type UserJob struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	Institution struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Links struct {
			Self   string `json:"self"`
			Source string `json:"source"`
		} `json:"links"`
	} `json:"institution"`
	Steps []struct {
		Title  string `json:"title"`
		Status string `json:"status"`
		Result struct {
			Code    string `json:"code"`
			Details string `json:"details"`
			Title   string `json:"title"`
			Type    string `json:"type"`
			URL     string `json:"url"`
		} `json:"result"`
	} `json:"steps"`
	Links struct {
		Self   string `json:"self"`
		Source string `json:"source"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) UserJobs(ctx context.Context, userID string) (UserJobList, error) {
	jobList, err := c.userJobs(ctx, userID)
	if err != nil && !IsUnauthorizedErr(err) {
		return jobList, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return UserJobList{}, err
	}
	return c.userJobs(ctx, userID)
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) userJobs(ctx context.Context, userID string) (UserJobList, error) {
	callURL, err := url.JoinPath(baseURL, "users", userID, "jobs")
	if err != nil {
		return UserJobList{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return UserJobList{}, err
	}

	var list UserJobList
	return list, json.Unmarshal(data, &list)
}
