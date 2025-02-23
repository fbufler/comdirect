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

	req.Header.Add("Accept", "application/json")

	_, _, err = c.doAuthenticatedRequest(req, token, http.StatusNoContent)
	if err != nil {
		return err
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

	req.Header.Add("Accept", "application/json")

	creationTime := time.Now()
	resBody, _, err := c.doAuthenticatedRequest(req, token, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var authResponse AuthResponse

	if err := json.NewDecoder(resBody).Decode(&authResponse); err != nil {
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

	addXHTTPRequestInfoHeader(req, token.SessionGUID, token.RequestID)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", fmt.Sprintf("qSession=%s", token.SessionGUID))

	resBody, _, err := c.doAuthenticatedRequest(req, token, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var sessions []Session
	if err := json.NewDecoder(resBody).Decode(&sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (c *Client) ValidateSession(token *AuthToken, sessionID string) (string, error) {
	slog.Debug("Validating session")
	currentSession := Session{Identifier: sessionID}
	currentSession.Activated2FA = true
	currentSession.SessionTanActive = true

	url := fmt.Sprintf("%s/session/clients/user/v1/sessions/%s/validate", c.config.APIURL, currentSession.Identifier)
	payload, err := json.Marshal(currentSession)
	if err != nil {
		return "", err
	}
	body := strings.NewReader(string(payload))
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return "", err
	}

	addXHTTPRequestInfoHeader(req, token.SessionGUID, token.RequestID)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resBody, header, err := c.doAuthenticatedRequest(req, token, http.StatusCreated)
	if err != nil {
		return "", err
	}

	var session Session
	if err := json.NewDecoder(resBody).Decode(&session); err != nil {
		return "", err
	}

	if session.SessionTanActive {
		/*
			TODO:
				if (header.typ==="P_TAN") {
					var image = 'data:image/png;base64,';
					image += header.challenge;
					var template = `<img src={{{data}}}></img>`;
					pm.visualizer.set(template, {data: image});
				}
		*/
		challengeHeader := header.Get("x-once-authentication-info")
		if challengeHeader == "" {
			return "", fmt.Errorf("missing x-once-authentication-info header")
		}
		var challenge xOnceAuthenticationInfo
		if err := json.Unmarshal([]byte(challengeHeader), &challenge); err != nil {
			return "", err
		}

		if challenge.ChallengeID == "" {
			return "", fmt.Errorf("missing challenge id in x-once-authentication-info header")
		}

		return challenge.ChallengeID, nil
	}

	return "", fmt.Errorf("session tan not active")
}

func (c *Client) ActivateSession(token *AuthToken, sessionID string, challengeId string) (*Session, error) {
	slog.Debug("Activating session")
	currentSession := Session{Identifier: sessionID}
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

	var session Session
	if err := json.NewDecoder(resBody).Decode(&session); err != nil {
		return nil, err
	}
	if !session.SessionTanActive {
		return nil, fmt.Errorf("session tan not active")
	}

	return &session, nil
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
