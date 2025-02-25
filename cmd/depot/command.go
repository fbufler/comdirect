package depot

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fbufler/comdirect/config"
	"github.com/fbufler/comdirect/internal/convert"
	"github.com/fbufler/comdirect/internal/flows"
	"github.com/fbufler/comdirect/pkg/comdirect"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "depot",
		Short: "Retrieve Depot Information",
	}

	cmd.PersistentFlags().StringP("output", "o", "json", "Output format (json, yaml)")
	cmd.AddCommand(depotsCmd)
	cmd.AddCommand(depotPositionCmd)
	cmd.AddCommand(depotPositionsCmd)
	cmd.AddCommand(depotTransactionsCmd)

	return cmd
}

var depotsCmd = &cobra.Command{
	Use:   "depots",
	Short: "Retrieve Depots",
	Run:   depots,
}

func depots(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	data, err := flows.Depots(cfg)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	handleOutput(cmd, data)
}

var depotPositionCmd = &cobra.Command{
	Use:   "position <depot-id> <position-id>",
	Short: "Retrieve Depot Position",
	Args:  cobra.ExactArgs(2),
	Run:   depotPosition,
}

func depotPosition(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	depotID := args[0]
	positionID := args[1]
	includeInstrument := cmd.Flag("include-instrument").Changed
	data, err := flows.DepotPosition(cfg, depotID, positionID, includeInstrument)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	handleOutput(cmd, data)
}

var depotPositionsCmd = &cobra.Command{
	Use:   "positions <depot-id>",
	Short: "Retrieve Depot Positions",
	Args:  cobra.ExactArgs(1),
	Run:   depotPositions,
}

func depotPositions(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	depotID := args[0]
	includeInstrument := cmd.Flag("include-instrument").Changed
	excludeDepot := cmd.Flag("exclude-depot").Changed
	data, err := flows.DepotPositions(cfg, depotID, includeInstrument, excludeDepot)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	handleOutput(cmd, data)
}

var depotTransactionsCmd = &cobra.Command{
	Use:   "transactions <depot-id>",
	Short: "Retrieve Depot Transactions",
	Args:  cobra.ExactArgs(1),
	Run:   depotTransactions,
}

func depotTransactions(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	depotID := args[0]
	wkn := cmd.Flag("wkn").Value.String()
	isin := cmd.Flag("isin").Value.String()
	instrumentID := cmd.Flag("instrument-id").Value.String()
	bookingStatus := comdirect.BookingStatus(cmd.Flag("booking-status").Value.String())
	maxBookingDateInput := cmd.Flag("max-booking-date").Value.String()
	var maxBookingDate time.Time
	if maxBookingDateInput != "" {
		var err error
		maxBookingDate, err = convert.TimeStringToTime(maxBookingDateInput)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	}
	data, err := flows.DepotTransactions(cfg, depotID, wkn, isin, instrumentID, bookingStatus, maxBookingDate)
	if err != nil {
		cmd.PrintErrln(err)
		return
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
	depotPositionCmd.Flags().Bool("include-instrument", false, "Include Instrument Information")

	depotPositionsCmd.Flags().Bool("include-instrument", false, "Include Instrument Information")
	depotPositionsCmd.Flags().Bool("exclude-depot", false, "Exclude Depot Information")

	depotTransactionsCmd.Flags().String("wkn", "", "WKN")
	depotTransactionsCmd.Flags().String("isin", "", "ISIN")
	depotTransactionsCmd.Flags().String("instrument-id", "", "Instrument ID")
	depotTransactionsCmd.Flags().String("booking-status", "", "Booking Status")
	depotTransactionsCmd.Flags().String("max-booking-date", "", "Max Booking Date e.g. 2006-01-02, 2006/01/02, 01/02/2006, 02.01.2006, 02.01.06")
}
