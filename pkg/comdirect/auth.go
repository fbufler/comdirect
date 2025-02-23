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

func (c *Client) Authenticate(twoFaHandler func(tanHeader TANHeader) error) (*AuthToken, error) {
	token, err := c.newInitialToken()
	if err != nil {
		return nil, err
	}

	sessions, err := c.sessions(token)
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("no session found")
	}

	// TODO: handle multiple sessions
	guessedSession := sessions[0]

	sessionGUID := token.SessionGUID
	challengeID, err := c.validateSession(token, guessedSession.Identifier)
	if err != nil {
		return nil, err
	}

	err = twoFaHandler(*challengeID)
	if err != nil {
		return nil, err
	}

	_, err = c.activateSession(token, sessionGUID, challengeID.Id)
	if err != nil {
		return nil, err
	}

	secondaryToken, err := c.newSecondaryToken(token)
	if err != nil {
		return nil, err
	}

	return secondaryToken, nil
}

func (c *Client) RefreshToken(token *AuthToken) (*AuthToken, error) {
	slog.Debug("Refreshing token")
	payload := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=refresh_token&refresh_token=%s", c.config.ClientID, c.config.ClientSecret, token.RefreshToken)
	body := strings.NewReader(payload)
	req, err := http.NewRequest(http.MethodPost, c.config.TokenURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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

	var authResponse authResponse

	if err := json.NewDecoder(res.Body).Decode(&authResponse); err != nil {
		return nil, err
	}

	if authResponse.AccessToken == "" {
		return nil, fmt.Errorf("missing access token in response")
	}

	return &AuthToken{
		AccessToken:  authResponse.AccessToken,
		ExpiresIn:    authResponse.ExpiresIn,
		RefreshToken: authResponse.RefreshToken,
		CreationTime: creationTime,
		Scope:        authResponse.Scope,
		SessionGUID:  token.SessionGUID,
		RequestID:    token.RequestID,
	}, nil
}

func (c *Client) RevokeToken(token *AuthToken) error {
	slog.Debug("Revoking token")
	c.ensureValidToken(token)
	req, err := http.NewRequest(http.MethodDelete, c.config.RevokeTokenURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")

	_, _, err = c.doAuthenticatedRequest(req, token, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) newInitialToken() (*AuthToken, error) {
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

	var authResponse authResponse
	if err := json.NewDecoder(res.Body).Decode(&authResponse); err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken:  authResponse.AccessToken,
		ExpiresIn:    authResponse.ExpiresIn,
		RefreshToken: authResponse.RefreshToken,
		CreationTime: creationTime,
		Scope:        authResponse.Scope,
		SessionGUID:  sessionID,
		RequestID:    requestID(),
	}, nil
}

func (c *Client) sessions(token *AuthToken) ([]session, error) {
	slog.Debug("Checking session status")
	c.ensureValidToken(token)
	url := fmt.Sprintf("%s/session/clients/user/v1/sessions", c.config.APIURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, token.SessionGUID, token.RequestID)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", fmt.Sprintf("qSession=%s", token.SessionGUID))

	resBody, _, err := c.doAuthenticatedRequest(req, token, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var sessions []session
	if err := json.NewDecoder(resBody).Decode(&sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (c *Client) validateSession(token *AuthToken, sessionID string) (*TANHeader, error) {
	slog.Debug("Validating session")
	currentSession := session{Identifier: sessionID}
	currentSession.Activated2FA = true
	currentSession.SessionTanActive = true

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

	addXHTTPRequestInfoHeader(req, token.SessionGUID, token.RequestID)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resBody, header, err := c.doAuthenticatedRequest(req, token, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	var session session
	if err := json.NewDecoder(resBody).Decode(&session); err != nil {
		return nil, err
	}

	if session.SessionTanActive {
		xoaiHeader := header.Get("x-once-authentication-info")
		if xoaiHeader == "" {
			return nil, fmt.Errorf("missing x-once-authentication-info header")
		}
		var tanHeader TANHeader
		if err := json.Unmarshal([]byte(xoaiHeader), &tanHeader); err != nil {
			return nil, err
		}

		if tanHeader.Id == "" {
			return nil, fmt.Errorf("missing challenge id in x-once-authentication-info header")
		}

		return &tanHeader, nil
	}

	return nil, fmt.Errorf("session tan not active")
}

func (c *Client) activateSession(token *AuthToken, sessionID string, challengeId string) (*session, error) {
	slog.Debug("Activating session")
	currentSession := session{Identifier: sessionID}
	currentSession.Activated2FA = true
	currentSession.SessionTanActive = true

	url := fmt.Sprintf("%s/session/clients/user/v1/sessions/%s", c.config.APIURL, currentSession.Identifier)
	payload, err := json.Marshal(currentSession)
	if err != nil {
		return nil, err
	}
	body := strings.NewReader(string(payload))
	req, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, token.SessionGUID, token.RequestID)
	addXOnceAuthenticationInfoHeader(req, challengeId)
	addXOnceAuthenticationHeader(req, "000000")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resBody, _, err := c.doAuthenticatedRequest(req, token, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var session session
	if err := json.NewDecoder(resBody).Decode(&session); err != nil {
		return nil, err
	}
	if !session.SessionTanActive {
		return nil, fmt.Errorf("session tan not active")
	}

	return &session, nil
}

func (c *Client) newSecondaryToken(token *AuthToken) (*AuthToken, error) {
	slog.Debug("Getting secondary token")

	payload := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=cd_secondary&token=%s", c.config.ClientID, c.config.ClientSecret, token.AccessToken)
	body := strings.NewReader(payload)

	req, err := http.NewRequest(http.MethodPost, c.config.TokenURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, handleRequestError(res)
	}

	defer res.Body.Close()

	var authResponse authResponse
	if err := json.NewDecoder(res.Body).Decode(&authResponse); err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken:  authResponse.AccessToken,
		ExpiresIn:    authResponse.ExpiresIn,
		RefreshToken: authResponse.RefreshToken,
		Scope:        authResponse.Scope,
		CreationTime: time.Now(),
		SessionGUID:  token.SessionGUID,
		RequestID:    token.RequestID,
	}, nil
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
