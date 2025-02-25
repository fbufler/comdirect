package main

import (
	"github.com/fbufler/comdirect/cmd/account"
	"github.com/fbufler/comdirect/cmd/depot"
	"github.com/fbufler/comdirect/cmd/e2e"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "comdirect",
	Short: "comdirect is a Go client for the comdirect API",
}

func init() {
	rootCmd.AddCommand(e2e.Command())
	rootCmd.AddCommand(account.Command())
	rootCmd.AddCommand(depot.Command())
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func main() {
	rootCmd.Execute()
}
