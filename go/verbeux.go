package verbeux

import (
	"net/http"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type Option func(*Client)

func NewClient(opts ...Option) *Client {
	c := &Client{
		apiKey:     "",
		baseURL:    defaultURL,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithHttpClient sets the http client of the client
func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithApiKey sets the apikey of the client
func WithApiKey(apiKey string) Option {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// WithBaseUrl sets the base url of the client
func WithBaseUrl(baseUrl string) Option {
	return func(c *Client) {
		c.baseURL = baseUrl
	}
}
