package comdirect

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (c *Client) NewToken() (*AuthToken, error) {
	slog.Debug("Getting token")
	sessionID := uuid.New().String()
	payload := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=password&username=%s&password=%s", c.config.ClientID, c.config.ClientSecret, c.config.Zugangsnummer, c.config.Pin)
	body := strings.NewReader(payload)

	req, err := http.NewRequest(http.MethodPost, c.config.TokenURL, body)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", fmt.Sprintf("qSession=%s", sessionID))

	creationTime := time.Now()
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, handleRequestError(res)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&authResponse); err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken:  authResponse.AccessToken,
		ExpiresIn:    authResponse.ExpiresIn,
		RefreshToken: authResponse.RefreshToken,
		CreationTime: creationTime,
		SessionGUID:  sessionID,
		RequestID:    requestID(),
	}, nil
}

func (c *Client) RevokeToken(token *AuthToken) error {
	slog.Debug("Revoking token")
	c.ensureValidToken(token)
	req, err := http.NewRequest(http.MethodDelete, c.config.RevokeTokenURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Accept", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 204 {
		return handleRequestError(res)
	}

	return nil
}

func (c *Client) RefreshToken(token *AuthToken) (*AuthToken, error) {
	slog.Debug("Refreshing token")
	payload := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=refresh_token&refresh_token=%s", c.config.ClientID, c.config.ClientSecret, token.RefreshToken)
	body := strings.NewReader(payload)
	req, err := http.NewRequest(http.MethodGet, c.config.TokenURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Accept", "application/json")

	creationTime := time.Now()
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, handleRequestError(res)
	}

	var authResponse AuthResponse

	if err := json.NewDecoder(res.Body).Decode(&authResponse); err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken:  authResponse.AccessToken,
		ExpiresIn:    authResponse.ExpiresIn,
		RefreshToken: authResponse.RefreshToken,
		CreationTime: creationTime,
	}, nil
}

func (c *Client) Sessions(token *AuthToken) ([]Session, error) {
	slog.Debug("Checking session status")
	c.ensureValidToken(token)
	url := fmt.Sprintf("%s/session/clients/user/v1/sessions", c.config.APIURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-http-request-info", fmt.Sprintf("{\"clientRequestId\":{\"sessionId\":\"%s\",\"requestId\":\"%s\"}}", token.SessionGUID, token.RequestID))
	req.Header.Add("Cookie", fmt.Sprintf("qSession=%s", token.SessionGUID))

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, handleRequestError(res)
	}

	var sessions []Session
	if err := json.NewDecoder(res.Body).Decode(&sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (c *Client) ValidateSession(token *AuthToken, currentSession *Session) (interface{}, error) {
	slog.Debug("Validating session")
	c.ensureValidToken(token)
	url := fmt.Sprintf("%s/session/clients/user/v1/sessions/%s/validate", c.config.APIURL, currentSession.Identifier)
	payload, err := json.Marshal(currentSession)
	if err != nil {
		return nil, err
	}
	body := strings.NewReader(string(payload))
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-http-request-info", fmt.Sprintf("{\"clientRequestId\":{\"sessionId\":\"%s\",\"requestId\":\"%s\"}}", token.SessionGUID, token.RequestID))

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, handleRequestError(res)
	}

	var session interface{}
	if err := json.NewDecoder(res.Body).Decode(&session); err != nil {
		return nil, err
	}

	return session, nil
}

func (c *Client) ensureValidToken(token *AuthToken) error {
	if time.Now().After(token.CreationTime.Add(time.Duration(token.ExpiresIn) * time.Second)) {
		slog.Debug("Token expired, refreshing")
		newToken, err := c.RefreshToken(token)
		if err != nil {
			slog.Debug("Token refresh failed")
			return err
		}
		token = newToken
	}
	return nil
}

func requestID() string {
	time := time.Now()
	return fmt.Sprintf("%d", time.UnixMilli())[10:]
}
