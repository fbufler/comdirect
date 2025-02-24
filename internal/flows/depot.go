package flows

import (
	"encoding/json"
	"time"

	"github.com/fbufler/comdirect/config"
	"github.com/fbufler/comdirect/pkg/comdirect"
)

func Depots(cfg *config.Config) (string, error) {
	client, token, err := bootstrap(cfg)
	if err != nil {
		return "", err
	}

	depots, err := client.Depots(token)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(depots)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func DepotPosition(cfg *config.Config, depotID, positionID string, includeInstrument bool) (string, error) {
	options := &comdirect.DepotPositionOptions{
		IncludeInstrument: includeInstrument,
	}
	client, token, err := bootstrap(cfg)
	if err != nil {
		return "", err
	}

	depot, err := client.DepotPosition(token, depotID, positionID, options)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(depot)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func DepotPositions(cfg *config.Config, depotID string, includeInstrument, excludeDepot bool) (string, error) {
	options := &comdirect.DepotOptions{
		IncludeInstrument: includeInstrument,
		ExcludeDepot:      excludeDepot,
	}

	client, token, err := bootstrap(cfg)
	if err != nil {
		return "", err
	}

	depot, err := client.DepotPositions(token, depotID, options)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(depot)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func DepotTransactions(cfg *config.Config, depotID string, wkn, isin, instrumentId string, bookingStatus comdirect.BookingStatus, maxBookingDate time.Time) (string, error) {
	options := &comdirect.DepotTransactionOptions{
		WKN:            wkn,
		ISIN:           isin,
		InstrumentId:   instrumentId,
		BookingStatus:  bookingStatus,
		MaxBookingDate: maxBookingDate,
	}
	client, token, err := bootstrap(cfg)
	if err != nil {
		return "", err
	}

	depotTransactions, err := client.DepotTransactions(token, depotID, options)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(depotTransactions)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
