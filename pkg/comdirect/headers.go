package comdirect

import (
	"fmt"
	"net/http"
)

const (
	authorizationHeader           = "Authorization"
	xHTTPRequestInfoHeader        = "x-http-request-info"
	xOnceAuthenticationInfoHeader = "x-once-authentication-info"
	XOnceAuthenticationHeader     = "x-once-authentication"
)

type xOnceAuthenticationInfo struct {
	ChallengeID  string `json:"id"`
	ChallengeTyp string `json:"typ"`
}

func addAuthorizationHeader(req *http.Request, token *AuthToken) {
	req.Header.Add(authorizationHeader, fmt.Sprintf("Bearer %s", token.AccessToken))
}

func addXHTTPRequestInfoHeader(req *http.Request, sessionGUID, requestID string) {
	req.Header.Add(xHTTPRequestInfoHeader, fmt.Sprintf("{\"clientRequestId\":{\"sessionId\":\"%s\",\"requestId\":\"%s\"}}", sessionGUID, requestID))
}

func addXOnceAuthenticationInfoHeader(req *http.Request, challengeID string) {
	req.Header.Add(xOnceAuthenticationInfoHeader, fmt.Sprintf("{\"id\":\"%s\"}", challengeID))
}

func addXOnceAuthenticationHeader(req *http.Request, tan string) {
	req.Header.Add(XOnceAuthenticationHeader, tan)
}
