package main

import (
	"github.com/fbufler/comdirect/cmd/e2e"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "comdirect",
	Short: "comdirect is a Go client for the comdirect API",
}

func init() {
	rootCmd.AddCommand(e2e.Command())
}

func main() {
	rootCmd.Execute()
}
