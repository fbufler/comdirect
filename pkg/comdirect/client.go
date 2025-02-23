package comdirect

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

const requestLimitPerSecond = 10

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
	config                Config
	client                *http.Client
	tokenStore            map[string]*AuthToken
	requestMonitor        map[time.Time]int
	requestLimitPerSecond int
}

func NewClient(config Config) *Client {
	return &Client{config: config, client: &http.Client{}, tokenStore: make(map[string]*AuthToken), requestMonitor: make(map[time.Time]int), requestLimitPerSecond: requestLimitPerSecond}
}

// AutoRefreshToken checks if any token is about to expire and refreshes it
// This function is blocking so it should be run in a goroutine
// As tokens are passed by reference, your token will be updated automatically if you have a reference to it
// This function is rather ment for a long idle time or a long running application
func (c *Client) AutoRefreshToken(ctx context.Context, expirationThreshold time.Duration) {
	expiringTokens := make(chan *AuthToken)
	go c.tokenRefresher(ctx, expiringTokens)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			{
				for _, token := range c.tokenStore {
					if token.WillExpireIn(expirationThreshold) {
						expiringTokens <- token
					}
				}
			}
		}
	}
}

func (c *Client) tokenRefresher(ctx context.Context, expiringTokens chan *AuthToken) {
	for {
		select {
		case <-ctx.Done():
			return
		case token := <-expiringTokens:
			{
				refreshedToken, err := c.RefreshToken(token)
				if err != nil {
					if err.Error() == LockedTokenError {
						slog.Warn("Token is locked, skipping refresh")
					} else {
						slog.Error(err.Error())
					}
					continue
				}
				token = refreshedToken
			}
		}
	}
}
