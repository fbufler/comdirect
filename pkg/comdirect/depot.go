package comdirect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DepotsOptions struct {
	PagingFirst int
}

func (o *DepotsOptions) queryParams() []string {
	queryParams := []string{}
	if o.PagingFirst > 0 {
		queryParams = append(queryParams, fmt.Sprintf("paging-first=%d", o.PagingFirst))
	}
	return queryParams
}

func (c *Client) Depots(authToken *AuthToken, options *DepotsOptions) (*Depots, error) {
	url := fmt.Sprintf("%s/brokerage/clients/user/v3/depots", c.config.APIURL)

	if options != nil {
		addQueryParams(url, options)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, authToken.SessionGUID, authToken.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.authenticatedRequest(req, authToken, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var depots Depots
	if err := json.NewDecoder(resBody).Decode(&depots); err != nil {
		return nil, err
	}

	return &depots, nil
}

func (c *Client) PaginatedDepots(authToken *AuthToken, amount int) (*Depots, error) {
	options := &DepotsOptions{
		PagingFirst: 0,
	}
	firstDepots, err := c.Depots(authToken, options)
	if err != nil {
		return nil, err
	}

	var data []Depots
	data = append(data, *firstDepots)
	for len(data) < amount/globalPageSize {
		options.PagingFirst = len(data) * globalPageSize
		page, err := c.Depots(authToken, options)
		if err != nil {
			return nil, err
		}

		if len(page.Values) == 0 {
			break
		}

		data = append(data, *page)
	}

	depots := flattenDepots(data)
	return &depots, nil
}

func flattenDepots(depots []Depots) Depots {
	values := []Depot{}
	for _, d := range depots {
		values = append(values, d.Values...)
	}
	return Depots{
		Paging: Paging{
			Index:   0,
			Matches: len(values),
		},
		Values: values,
	}
}

type DepotPosistionsOptions struct {
	IncludeInstrument bool
	ExcludeDepot      bool
	PagingFirst       int
}

func (o *DepotPosistionsOptions) queryParams() []string {
	queryParams := []string{}
	if o.IncludeInstrument {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", includeProperty, "instrument"))
	}
	if o.ExcludeDepot {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", excludeProperty, "depot"))
	}
	if o.PagingFirst > 0 {
		queryParams = append(queryParams, fmt.Sprintf("paging-first=%d", o.PagingFirst))
	}
	return queryParams
}

// DepotPositions returns the positions of a depot.
// For more information see https://www.comdirect.de/cms/media/comdirect_REST_API_Dokumentation.pdf
func (c *Client) DepotPositions(authToken *AuthToken, depotID string, options *DepotPosistionsOptions) (*DepotPositions, error) {
	url := fmt.Sprintf("%s/brokerage/v3/depots/%s/positions", c.config.APIURL, depotID)

	if options != nil {
		addQueryParams(url, options)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, authToken.SessionGUID, authToken.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.authenticatedRequest(req, authToken, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var depotPositions DepotPositions
	if err := json.NewDecoder(resBody).Decode(&depotPositions); err != nil {
		return nil, err
	}

	return &depotPositions, nil
}

func (c *Client) PaginatedDepotPositions(authToken *AuthToken, depotID string, amount int, options *DepotPosistionsOptions) (*DepotPositions, error) {
	if options == nil {
		options = &DepotPosistionsOptions{}
	}
	options.PagingFirst = 0
	firstPositions, err := c.DepotPositions(authToken, depotID, options)
	if err != nil {
		return nil, err
	}

	var data []DepotPositions
	data = append(data, *firstPositions)
	for len(data) < amount/globalPageSize {
		options.PagingFirst = len(data) * globalPageSize
		page, err := c.DepotPositions(authToken, depotID, options)
		if err != nil {
			return nil, err
		}

		if len(page.Values) == 0 {
			break
		}

		data = append(data, *page)
	}

	positions := flattenDepotPositions(data)
	return &positions, nil
}

func flattenDepotPositions(positions []DepotPositions) DepotPositions {
	values := []DepotPosition{}
	for _, p := range positions {
		values = append(values, p.Values...)
	}
	return DepotPositions{
		Paging: Paging{
			Index:   0,
			Matches: len(values),
		},
		Values: values,
	}
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
		addQueryParams(url, options)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, authToken.SessionGUID, authToken.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.authenticatedRequest(req, authToken, http.StatusOK)
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
	PagingFirst    int
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
	if dto.PagingFirst > 0 {
		queryParams = append(queryParams, fmt.Sprintf("paging-first=%d", dto.PagingFirst))
	}
	return queryParams
}

func (c *Client) DepotTransactions(authToken *AuthToken, depotID string, options *DepotTransactionOptions) (*DepotTransactions, error) {
	url := fmt.Sprintf("%s/brokerage/v3/depots/%s/transactions", c.config.APIURL, depotID)

	if options != nil {
		addQueryParams(url, options)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addXHTTPRequestInfoHeader(req, authToken.SessionGUID, authToken.RequestID)
	req.Header.Add("Accept", "application/json")

	resBody, _, err := c.authenticatedRequest(req, authToken, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var depotTransactions DepotTransactions
	if err := json.NewDecoder(resBody).Decode(&depotTransactions); err != nil {
		return nil, err
	}

	return &depotTransactions, nil
}

func (c *Client) PaginatedDepotTransactions(authToken *AuthToken, depotID string, amount int, options *DepotTransactionOptions) (*DepotTransactions, error) {
	if options == nil {
		options = &DepotTransactionOptions{}
	}
	options.PagingFirst = 0
	firstPage, err := c.DepotTransactions(authToken, depotID, options)
	if err != nil {
		return nil, err
	}

	transactions := []*DepotTransactions{firstPage}
	for len(transactions) < amount/globalPageSize {
		options.PagingFirst = len(transactions) * globalPageSize
		page, err := c.DepotTransactions(authToken, depotID, options)
		if err != nil {
			return nil, err
		}
		if len(page.Values) == 0 {
			break
		}
		transactions = append(transactions, page)
	}

	return &DepotTransactions{
		Paging: Paging{
			Index:   0,
			Matches: len(transactions) * globalPageSize,
		},
		Values: flattenDepotTransactions(transactions),
	}, nil
}

func flattenDepotTransactions(transactions []*DepotTransactions) []DepotTransaction {
	var flattened []DepotTransaction
	for _, t := range transactions {
		flattened = append(flattened, t.Values...)
	}
	return flattened
}
