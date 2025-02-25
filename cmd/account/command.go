package account

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/fbufler/comdirect/config"
	"github.com/fbufler/comdirect/internal/convert"
	"github.com/fbufler/comdirect/internal/flows"
	"github.com/fbufler/comdirect/pkg/comdirect"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Retrieve Bank Account Information",
	}
	cmd.PersistentFlags().StringP("output", "o", "json", "Output format (json, yaml)")
	cmd.AddCommand(balancesCmd)
	cmd.AddCommand(balanceCmd)
	cmd.AddCommand(transactionsCmd)
	return cmd
}

var balancesCmd = &cobra.Command{
	Use:   "balances",
	Short: "Retrieve Account Balances",
	Run:   balances,
}

func balances(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	excludeAccount := cmd.Flag("exclude-account").Changed
	data, err := flows.AccountBalances(cfg, excludeAccount)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	handleOutput(cmd, data)
}

var balanceCmd = &cobra.Command{
	Use:   "balance <account-id>",
	Short: "Retrieve Account Balance",
	Args:  cobra.ExactArgs(1),
	Run:   balance,
}

func balance(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	accountID := args[0]
	data, err := flows.AccountBalance(cfg, accountID)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	handleOutput(cmd, data)
}

var transactionsCmd = &cobra.Command{
	Use:   "transactions <account-id>",
	Short: "Retrieve Account Transactions",
	Args:  cobra.ExactArgs(1),
	Run:   transactions,
}

func transactions(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	accountID := args[0]
	transactionState := comdirect.TransactionState(cmd.Flag("state").Value.String())
	includeAccount := cmd.Flag("include-account").Changed
	countInput := cmd.Flag("count").Value.String()

	var data string
	var err error
	if countInput != "" {
		count, err := strconv.Atoi(countInput)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		data, err = flows.PaginatedAccountTransactions(cfg, accountID, count, includeAccount)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	} else {
		data, err = flows.AccountTransactions(cfg, accountID, transactionState, includeAccount)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	}
	handleOutput(cmd, data)
}

func handleOutput(cmd *cobra.Command, data string) {
	output := cmd.Flag("output").Value.String()
	slog.Info(fmt.Sprintf("Output format: %s", output))
	switch output {
	case "json":
		json, err := convert.JSONToReadableJSON(data)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		cmd.Println(json)
	case "yaml":
		yaml, err := convert.JSONToYAML(data)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		cmd.Println(yaml)
	default:
		cmd.PrintErrln("Unsupported output format")
	}
}

func init() {
	balancesCmd.Flags().BoolP("exclude-account", "e", false, "Exclude Account")
	transactionsCmd.Flags().StringP("state", "s", string(comdirect.TransactionStateBoth), "Transaction State (BOTH, BOOKED, NOTBOOKED)")
	transactionsCmd.Flags().BoolP("include-account", "i", false, "Include Account")
	transactionsCmd.Flags().StringP("count", "c", "", "Amount of Transactions, by default 20")
}
