package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/fbufler/comdirect/pkg/comdirect"
)

// e2e test
func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	// load from env
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	zugangsnummer := os.Getenv("ZUGANGSNUMMER")
	pin := os.Getenv("PIN")
	apiURL := os.Getenv("API_URL")
	tokenURL := os.Getenv("TOKEN_URL")
	revokeTokenURL := os.Getenv("REVOKE_TOKEN_URL")

	// create config
	config := comdirect.Config{
		APIURL:         apiURL,
		TokenURL:       tokenURL,
		RevokeTokenURL: revokeTokenURL,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Zugangsnummer:  zugangsnummer,
		Pin:            pin,
	}

	fmt.Println(config)

	// create client
	client := comdirect.NewClient(config)

	// Authenticate
	token, err := client.Authenticate(twoFaHandler)
	if err != nil {
		panic(err)
	}

	// Refresh token
	token, err = client.RefreshToken(token)
	if err != nil {
		panic(err)
	}

	// Get account balances
	accountBalances, err := client.AccountBalances(token, nil)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Account balances: %v", accountBalances))

	// Get account balance
	relevantAccountID := accountBalances.Values[1].AccountID
	accountBalance, err := client.AccountBalance(token, relevantAccountID)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Account balance: %v", accountBalance))

	// Get transactions
	transactions, err := client.AccountTransactions(token, relevantAccountID, comdirect.TransactionStateBooked, nil)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Transactions: %v", transactions))

	// Revoke token
	err = client.RevokeToken(token)
	if err != nil {
		panic(err)
	}
}

func twoFaHandler(tanHeader comdirect.TANHeader) error {
	slog.Info("Please verify the TAN")
	slog.Info(fmt.Sprintf("TAN - id: %s - typ: %s", tanHeader.Id, tanHeader.Typ))
	// wait for user input

	slog.Info("Press enter to continue")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	slog.Info("Continuing")
	return nil
}
