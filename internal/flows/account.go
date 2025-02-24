package flows

import (
	"encoding/json"

	"github.com/fbufler/comdirect/config"
	"github.com/fbufler/comdirect/pkg/comdirect"
)

func AccountBalances(cfg *config.Config, excludeAccount bool) (string, error) {
	options := &comdirect.AccountBalancesOptions{
		ExludeAccount: excludeAccount,
	}
	client, token, err := bootstrap(cfg)
	if err != nil {
		return "", err
	}

	accounts, err := client.AccountBalances(token, options)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(accounts)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func AccountBalance(cfg *config.Config, accountID string) (string, error) {
	client, token, err := bootstrap(cfg)
	if err != nil {
		return "", err
	}

	account, err := client.AccountBalance(token, accountID)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(account)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func AccountTransactions(cfg *config.Config, accountID string, transactionState comdirect.TransactionState, includeAccount bool) (string, error) {
	options := &comdirect.AccountTransactionOptions{
		IncludeAccount: includeAccount,
	}
	client, token, err := bootstrap(cfg)
	if err != nil {
		return "", err
	}

	transactions, err := client.AccountTransactions(token, accountID, transactionState, options)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(transactions)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
