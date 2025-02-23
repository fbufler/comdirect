package comdirect

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) AccountBalance(token *AuthToken) ([]interface{}, error) {
	url := fmt.Sprintf("%s/banking/clients/user/v2/accounts/balances", c.config.APIURL)
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

	var accountBalances []interface{}
	if err := json.NewDecoder(resBody).Decode(&accountBalances); err != nil {
		return nil, err
	}

	return accountBalances, nil
}
