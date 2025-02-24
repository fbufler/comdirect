package comdirect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	includeProperty = "with-attr"
	excludeProperty = "without-attr"
)

type Projection struct {
	IncludeProperties []string
	ExcludeProperties []string
}

func (p *Projection) queryParams() []string {
	var params []string

	if len(p.IncludeProperties) > 0 {
		params = append(params, fmt.Sprintf("%s=%s", includeProperty, strings.Join(p.IncludeProperties, ",")))
	}

	if len(p.ExcludeProperties) > 0 {
		params = append(params, fmt.Sprintf("%s=%s", excludeProperty, strings.Join(p.ExcludeProperties, ",")))
	}

	return params
}

type AccountBalancesOptions struct {
	ExludeAccount bool
}

func (o *AccountBalancesOptions) queryParams() []string {
	queryParams := []string{}
	if o.ExludeAccount {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", excludeProperty, "account"))
	}
	return queryParams
}

// AccountBalances returns the balances of all accounts of the user.
// For more information see https://www.comdirect.de/cms/media/comdirect_REST_API_Dokumentation.pdf
func (c *Client) AccountBalances(token *AuthToken, options *AccountBalancesOptions) (*AccountBalances, error) {
	url := fmt.Sprintf("%s/banking/clients/user/v2/accounts/balances", c.config.APIURL)
	if options != nil {
		queryParams := options.queryParams()
		if len(queryParams) > 0 {
			url += fmt.Sprintf("?%s", strings.Join(queryParams, "&"))
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, token.SessionGUID, token.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.doAuthenticatedRequest(req, token, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var accountBalances AccountBalances
	if err := json.NewDecoder(resBody).Decode(&accountBalances); err != nil {
		return nil, err
	}

	return &accountBalances, nil
}

// AccountBalance returns the balance of a specific account.
// For more information see https://www.comdirect.de/cms/media/comdirect_REST_API_Dokumentation.pdf
func (c *Client) AccountBalance(token *AuthToken, accountID string) (*AccountBalance, error) {
	url := fmt.Sprintf("%s/banking/v2/accounts/%s/balances", c.config.APIURL, accountID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, token.SessionGUID, token.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.doAuthenticatedRequest(req, token, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var accountBalance AccountBalance
	if err := json.NewDecoder(resBody).Decode(&accountBalance); err != nil {
		return nil, err
	}

	return &accountBalance, nil
}

type TransactionState string

const (
	TransactionStateBoth      TransactionState = "BOTH"
	TransactionStateBooked    TransactionState = "BOOKED"
	TransactionStateNotBooked TransactionState = "NOTBOOKED"
)

type AccountTransactionOptions struct {
	IncludeAccount bool
}

func (o *AccountTransactionOptions) queryParams() []string {
	queryParams := []string{}
	if o.IncludeAccount {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", includeProperty, "account"))
	}
	return queryParams
}

// AccountTransactions returns the transactions of a specific account.
// For more information see https://www.comdirect.de/cms/media/comdirect_REST_API_Dokumentation.pdf
func (c *Client) AccountTransactions(token *AuthToken, accountID string, transactionState TransactionState, options *AccountTransactionOptions) (*AccountTransactions, error) {
	url := fmt.Sprintf("%s/banking/v1/accounts/%s/transactions", c.config.APIURL, accountID)

	queryParams := []string{
		fmt.Sprintf("transactionState=%s", transactionState),
	}
	if options != nil {
		queryParams = append(queryParams, options.queryParams()...)
	}

	url += fmt.Sprintf("?%s", strings.Join(queryParams, "&"))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, token.SessionGUID, token.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.doAuthenticatedRequest(req, token, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var accountTransactions AccountTransactions
	if err := json.NewDecoder(resBody).Decode(&accountTransactions); err != nil {
		return nil, err
	}

	return &accountTransactions, nil
}
