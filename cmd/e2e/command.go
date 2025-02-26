package e2e

import (
	"fmt"
	"log/slog"

	"github.com/fbufler/comdirect/config"
	"github.com/fbufler/comdirect/internal/flows"
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
	client, token, err := flows.Bootstrap(cfg)
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
	options := &comdirect.AccountTransactionOptions{
		TransactionState: comdirect.TransactionStateBoth,
	}
	transactions, err := client.AccountTransactions(token, relevantAccountID, options)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Transactions: %v", transactions))

	// Get paginated transactions
	slog.Info("Getting paginated transactions")
	paginatedTransactions, err := client.PaginatedAccountTransactions(token, relevantAccountID, 60, nil)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Paginated transactions: %v", paginatedTransactions))

	// Get depots
	slog.Info("Getting depots")
	depots, err := client.Depots(token, nil)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Depots: %v", depots))

	// Get paginated depots
	slog.Info("Getting paginated depots")
	paginatedDepots, err := client.PaginatedDepots(token, 60)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Paginated depots: %v", paginatedDepots))

	// Get depot positions
	slog.Info("Getting depot positions")
	depotPositions, err := client.DepotPositions(token, depots.Values[0].DepotID, nil)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Depot positions: %v", depotPositions))

	// Get Paginated depot positions
	slog.Info("Getting paginated depot positions")
	paginatedDepotPositions, err := client.PaginatedDepotPositions(token, depots.Values[0].DepotID, 60, nil)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Paginated depot positions: %v", paginatedDepotPositions))

	// Get depot position
	slog.Info("Getting depot position")
	depotPosition, err := client.DepotPosition(token, depotPositions.Values[0].DepotID, depotPositions.Values[0].PositionID, nil)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Depot position: %v", depotPosition))

	// Get depot transactions
	slog.Info("Getting depot transactions")
	depotTransactions, err := client.DepotTransactions(token, depots.Values[0].DepotID, nil)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Depot transactions: %v", depotTransactions))

	// Get paginated depot transactions
	slog.Info("Getting paginated depot transactions")
	paginatedDepotTransactions, err := client.PaginatedDepotTransactions(token, depots.Values[0].DepotID, 60, nil)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Paginated depot transactions: %v", paginatedDepotTransactions))

	// Revoke token
	slog.Info("Revoking token")
	err = client.RevokeToken(token)
	if err != nil {
		panic(err)
	}
}
