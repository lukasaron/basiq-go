package basiq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Job struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Steps   []struct {
		Title  string `json:"title"`
		Status string `json:"status"`
		Result struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"result"`
	} `json:"steps"`
	Links struct {
		Self   string `json:"self"`
		Source string `json:"source"`
	} `json:"links"`
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) Job(ctx context.Context, jobID string) (Job, error) {
	job, err := c.job(ctx, jobID)
	if err != nil && !IsUnauthorizedErr(err) {
		return job, err
	}
	if err = c.Authenticate(ctx); err != nil {
		return Job{}, err
	}
	return c.job(ctx, jobID)
}

// --------------------------------------------------------------------------------------------------------------------

func (c *Client) job(ctx context.Context, jobID string) (Job, error) {
	callURl, err := url.JoinPath(baseURL, "jobs", jobID)
	if err != nil {
		return Job{}, err
	}

	data, err := c.makeCall(ctx, http.MethodGet, callURl, nil)
	if err != nil {
		return Job{}, err
	}

	var job Job
	return job, json.Unmarshal(data, &job)
}
