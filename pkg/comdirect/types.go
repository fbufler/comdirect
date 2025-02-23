package comdirect

import "time"

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type AuthToken struct {
	AccessToken  string
	ExpiresIn    int
	RefreshToken string
	CreationTime time.Time
	SessionGUID  string
	RequestID    string
}

type Session struct {
	Identifier       string `json:"identifier"`
	SessionTanActive bool   `json:"sessionTanActive"`
	Activated2FA     bool   `json:"activated2FA"`
}
