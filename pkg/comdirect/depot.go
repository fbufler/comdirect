package comdirect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (c *Client) Depots(authToken *AuthToken) (*Depots, error) {
	url := fmt.Sprintf("%s/brokerage/clients/user/v3/depots", c.config.APIURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, authToken.SessionGUID, authToken.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.doAuthenticatedRequest(req, authToken, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var depots Depots
	if err := json.NewDecoder(resBody).Decode(&depots); err != nil {
		return nil, err
	}

	return &depots, nil
}

type DepotOptions struct {
	IncludeInstrument bool
	ExcludeDepot      bool
}

func (o *DepotOptions) queryParams() []string {
	queryParams := []string{}
	if o.IncludeInstrument {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", includeProperty, "instrument"))
	}
	if o.ExcludeDepot {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", excludeProperty, "depot"))
	}
	return queryParams
}

// DepotPositions returns the positions of a depot.
// For more information see https://www.comdirect.de/cms/media/comdirect_REST_API_Dokumentation.pdf
func (c *Client) DepotPositions(authToken *AuthToken, depotID string, options *DepotOptions) (*DepotPositions, error) {
	url := fmt.Sprintf("%s/brokerage/v3/depots/%s/positions", c.config.APIURL, depotID)

	if options != nil {
		queryParams := options.queryParams()
		if len(queryParams) > 0 {
			url += fmt.Sprintf("?%s", strings.Join(queryParams, "&"))
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, authToken.SessionGUID, authToken.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.doAuthenticatedRequest(req, authToken, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var depotPositions DepotPositions
	if err := json.NewDecoder(resBody).Decode(&depotPositions); err != nil {
		return nil, err
	}

	return &depotPositions, nil
}

type DepotPositionOptions struct {
	IncludeInstrument bool
}

func (o *DepotPositionOptions) queryParams() []string {
	queryParams := []string{}
	if o.IncludeInstrument {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", includeProperty, "instrument"))
	}
	return queryParams
}

// DepotPosition returns the position of a depot.
// For more information see https://www.comdirect.de/cms/media/comdirect_REST_API_Dokumentation.pdf
func (c *Client) DepotPosition(authToken *AuthToken, depotID string, positionID string, options *DepotPositionOptions) (*DepotPosition, error) {
	url := fmt.Sprintf("%s/brokerage/v3/depots/%s/positions/%s", c.config.APIURL, depotID, positionID)

	if options != nil {
		queryParams := options.queryParams()
		if len(queryParams) > 0 {
			url += fmt.Sprintf("?%s", strings.Join(queryParams, "&"))
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, authToken.SessionGUID, authToken.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.doAuthenticatedRequest(req, authToken, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var depotPosition DepotPosition
	if err := json.NewDecoder(resBody).Decode(&depotPosition); err != nil {
		return nil, err
	}

	return &depotPosition, nil
}

type BookingStatus string

const (
	BookingStatusBooked    BookingStatus = "BOOKED"
	BookingStatusNotBooked BookingStatus = "NOTBOOKED"
	BookingStatusBoth      BookingStatus = "BOTH"
)

type DepotTransactionOptions struct {
	WKN            string
	ISIN           string
	InstrumentId   string
	BookingStatus  BookingStatus
	MaxBookingDate time.Time
}

func (dto *DepotTransactionOptions) queryParams() []string {
	queryParams := []string{}
	if dto.WKN != "" {
		queryParams = append(queryParams, fmt.Sprintf("wkn=%s", dto.WKN))
	}
	if dto.ISIN != "" {
		queryParams = append(queryParams, fmt.Sprintf("isin=%s", dto.ISIN))
	}
	if dto.InstrumentId != "" {
		queryParams = append(queryParams, fmt.Sprintf("instrumentId=%s", dto.InstrumentId))
	}
	if dto.BookingStatus != "" {
		queryParams = append(queryParams, fmt.Sprintf("bookingStatus=%s", dto.BookingStatus))
	}
	if !dto.MaxBookingDate.IsZero() {
		queryParams = append(queryParams, fmt.Sprintf("maxBookingDate=%s", dto.MaxBookingDate.Format("2006-01-02")))
	}
	return queryParams
}

func (c *Client) DepotTransactions(authToken *AuthToken, depotID string, options *DepotTransactionOptions) (*DepotTransactions, error) {
	url := fmt.Sprintf("%s/brokerage/v3/depots/%s/transactions", c.config.APIURL, depotID)

	if options != nil {
		queryParams := options.queryParams()
		if len(queryParams) > 0 {
			url += fmt.Sprintf("?%s", strings.Join(queryParams, "&"))
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, authToken.SessionGUID, authToken.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.doAuthenticatedRequest(req, authToken, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var depotTransactions DepotTransactions
	if err := json.NewDecoder(resBody).Decode(&depotTransactions); err != nil {
		return nil, err
	}

	return &depotTransactions, nil
}
