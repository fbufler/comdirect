package e2e

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/fbufler/comdirect/config"
	"github.com/fbufler/comdirect/pkg/comdirect"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "e2e",
		Short: "Run the end to end test",
		Run: func(cmd *cobra.Command, args []string) {
			e2e()
		},
	}
	return cmd
}

// e2e test
func e2e() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	cfg := config.Get()
	// create config
	config := comdirect.Config{
		APIURL:         cfg.APIURL,
		TokenURL:       cfg.TokenURL,
		RevokeTokenURL: cfg.RevokeTokenURL,
		ClientID:       cfg.ClientID,
		ClientSecret:   cfg.ClientSecret,
		Zugangsnummer:  cfg.Zugangsnummer,
		Pin:            cfg.Pin,
	}

	fmt.Println(config)

	// create client
	client := comdirect.NewClient(config)

	// Authenticate
	slog.Info("Authenticating")
	token, err := client.Authenticate(twoFaHandler)
	if err != nil {
		panic(err)
	}

	// Refresh token
	slog.Info("Refreshing token")
	token, err = client.RefreshToken(token)
	if err != nil {
		panic(err)
	}

	// Get account balances
	slog.Info("Getting account balances")
	accountBalances, err := client.AccountBalances(token, nil)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Account balances: %v", accountBalances))

	// Get account balance
	slog.Info("Getting account balance")
	relevantAccountID := accountBalances.Values[1].AccountID
	accountBalance, err := client.AccountBalance(token, relevantAccountID)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Account balance: %v", accountBalance))

	// Get transactions
	slog.Info("Getting transactions")
	transactions, err := client.AccountTransactions(token, relevantAccountID, comdirect.TransactionStateBooked, nil)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Transactions: %v", transactions))

	// Revoke token
	slog.Info("Revoking token")
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
