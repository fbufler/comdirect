package comdirect

import (
	"net/http"
)

type Config struct {
	APIURL         string
	TokenURL       string
	RevokeTokenURL string
	ClientID       string
	ClientSecret   string
	Zugangsnummer  string
	Pin            string
}

type Client struct {
	config Config
	client *http.Client
}

func NewClient(config Config) *Client {
	return &Client{config: config, client: &http.Client{}}
}
