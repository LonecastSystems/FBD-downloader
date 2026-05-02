package fbd

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const homeURL string = "https://www.football-bet-data.com/"

type Client struct {
	httpClient     *http.Client
	baseURL        *url.URL
	enforceSession bool
}

func NewClient() (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}

	baseURL, err := url.Parse(homeURL)
	if err != nil {
		return nil, fmt.Errorf("parse home URL: %w", err)
	}

	client := &Client{
		baseURL: baseURL,
	}
	transport := &sessionCheckTransport{
		base:     http.DefaultTransport,
		jar:      jar,
		baseURL:  baseURL,
		enforcer: func() bool { return client.enforceSession },
	}

	client.httpClient = &http.Client{
		Timeout:   30 * time.Second,
		Jar:       jar,
		Transport: transport,
	}

	return client, nil
}

func (c *Client) ensureReady() error {
	if c == nil || c.httpClient == nil || c.baseURL == nil {
		return errors.New("client is not initialized")
	}
	return nil
}
