package comdirect

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func (c *Client) doAuthenticatedRequest(req *http.Request, token *AuthToken, expectedStatus int) (io.Reader, *http.Header, error) {
	c.ensureValidToken(token)
	addAuthorizationHeader(req, token)
	c.newAuthenticatedRequest()
	token.Lock()
	res, err := c.client.Do(req)
	token.Unlock()
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != expectedStatus {
		return nil, nil, handleRequestError(res)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	reader := bytes.NewReader(body)
	return reader, &res.Header, nil
}

func (c *Client) newAuthenticatedRequest() {
	currentTimeToSecond := time.Now().Truncate(time.Second)
	c.requestMonitor[currentTimeToSecond]++
	if c.requestMonitor[currentTimeToSecond] > c.requestLimitPerSecond {
		slog.Warn(fmt.Sprintf("Reaching request limit of %d for current second: %s", c.requestLimitPerSecond, currentTimeToSecond))
	}
}

func handleRequestError(res *http.Response) error {
	statusErr := fmt.Errorf("request failed with status code %d", res.StatusCode)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return statusErr
	}
	if len(body) != 0 {
		slog.Debug(string(body))
	}
	return statusErr
}

func requestID() string {
	time := time.Now()
	return fmt.Sprintf("%d", time.UnixMilli())[10:]
}
