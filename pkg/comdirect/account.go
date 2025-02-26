package comdirect

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	includeProperty = "with-attr"
	excludeProperty = "without-attr"
)

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
		url = addQueryParams(url, options)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, token.SessionGUID, token.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.authenticatedRequest(req, token, http.StatusOK)
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

	resBody, _, err := c.authenticatedRequest(req, token, http.StatusOK)
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
	IncludeAccount   bool
	PagingFirst      int
	TransactionState TransactionState
}

func (o *AccountTransactionOptions) queryParams() []string {
	queryParams := []string{}
	if o.IncludeAccount {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", includeProperty, "account"))
	}
	if o.PagingFirst > 0 {
		queryParams = append(queryParams, fmt.Sprintf("paging-first=%d", o.PagingFirst))
	}
	if o.TransactionState != "" {
		queryParams = append(queryParams, fmt.Sprintf("transactionState=%s", o.TransactionState))
	} else {
		queryParams = append(queryParams, fmt.Sprintf("transactionState=%s", TransactionStateBoth))
	}
	return queryParams
}

// AccountTransactions returns the transactions of a specific account.
// For more information see https://www.comdirect.de/cms/media/comdirect_REST_API_Dokumentation.pdf
func (c *Client) AccountTransactions(token *AuthToken, accountID string, options *AccountTransactionOptions) (*AccountTransactions, error) {
	url := fmt.Sprintf("%s/banking/v1/accounts/%s/transactions", c.config.APIURL, accountID)

	if options != nil {
		addQueryParams(url, options)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, token.SessionGUID, token.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.authenticatedRequest(req, token, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var accountTransactions AccountTransactions
	if err := json.NewDecoder(resBody).Decode(&accountTransactions); err != nil {
		return nil, err
	}

	return &accountTransactions, nil
}

func (c *Client) PaginatedAccountTransactions(token *AuthToken, accountID string, amount int, options *AccountTransactionOptions) (*AccountTransactions, error) {
	if options == nil {
		options = &AccountTransactionOptions{}
	}
	options.TransactionState = TransactionStateBooked
	firstPage, err := c.AccountTransactions(token, accountID, options)
	if err != nil {
		return nil, err
	}

	transactions := []*AccountTransactions{firstPage}
	for len(transactions) < amount/globalPageSize {
		options.PagingFirst = len(transactions) * globalPageSize
		page, err := c.AccountTransactions(token, accountID, options)
		if err != nil {
			return nil, err
		}
		if len(page.Values) == 0 {
			break
		}
		transactions = append(transactions, page)
	}

	return &AccountTransactions{
		Paging: Paging{
			Index:   0,
			Matches: len(transactions) * globalPageSize,
		},
		AggregatedTransactions: firstPage.AggregatedTransactions,
		Values:                 flattenTransactions(transactions),
	}, nil
}

func flattenTransactions(transactions []*AccountTransactions) []AccountTransaction {
	var flattened []AccountTransaction
	for _, page := range transactions {
		flattened = append(flattened, page.Values...)
	}
	return flattened
}
