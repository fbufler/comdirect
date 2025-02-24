package comdirect

import "time"

type AuthToken struct {
	AccessToken  string
	ExpiresIn    int
	RefreshToken string
	CreationTime time.Time
	Scope        string
	SessionGUID  string
	RequestID    string
	locked       bool
}

func (t *AuthToken) IsExpired() bool {
	return time.Since(t.CreationTime).Seconds() > float64(t.ExpiresIn)
}

func (t *AuthToken) WillExpireIn(seconds time.Duration) bool {
	return time.Since(t.CreationTime).Seconds()+seconds.Seconds() > float64(t.ExpiresIn)
}

func (t *AuthToken) Lock() {
	t.locked = true
}

func (t *AuthToken) Unlock() {
	t.locked = false
}

func (t *AuthToken) IsLocked() bool {
	return t.locked
}

type authResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	KDNR         string `json:"kdnr"`
	BPID         int    `json:"bpid"`
	KontaktId    int    `json:"kontaktId"`
}

type session struct {
	Identifier       string `json:"identifier"`
	SessionTanActive bool   `json:"sessionTanActive"`
	Activated2FA     bool   `json:"activated2FA"`
}

type TANHeader struct {
	Id        string `json:"id"`
	Typ       string `json:"typ"`
	Challenge string `json:"challenge"`
}

type Paging struct {
	Index   int `json:"index"`
	Matches int `json:"matches"`
}

type Balance struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type AccountBalances struct {
	Paging Paging           `json:"paging"`
	Values []AccountBalance `json:"values"`
}

type AccountBalance struct {
	Account                Account `json:"account"`
	AccountID              string  `json:"accountId"`
	Balance                Balance `json:"balance"`
	BalanceEUR             Balance `json:"balanceEUR"`
	AvailableCashAmount    Balance `json:"availableCashAmount"`
	AvailableCashAmountEUR Balance `json:"availableCashAmountEUR"`
}

type Account struct {
	AccountID        string      `json:"accountId"`
	AccountDisplayID string      `json:"accountDisplayId"`
	Currency         string      `json:"currency"`
	ClientID         string      `json:"clientId"`
	AccountType      AccountType `json:"accountType"`
	IBAN             string      `json:"iban"`
	BIC              string      `json:"bic"`
	CreditLimit      Balance     `json:"creditLimit"`
}

type AccountType struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

type AccountTransactions struct {
	Paging                 Paging                        `json:"paging"`
	AggregatedTransactions AggregatedAccountTransactions `json:"aggregated"`
	Values                 []AccountTransaction          `json:"values"`
}

type AggregatedAccountTransactions struct {
	Account                      Account `json:"account"`
	AccountID                    string  `json:"accountId"`
	BookingDateLatestTransaction string  `json:"bookingDateLatestTransaction"`
	ReferenceLatestTransaction   string  `json:"referenceLatestTransaction"`
	LatestTransactionIncluded    bool    `json:"latestTransactionIncluded"`
	PagingTimestamp              string  `json:"pagingTimestamp"`
}

type AccountTransaction struct {
	Reference             string                 `json:"reference"`
	BookingStatus         string                 `json:"bookingStatus"`
	BookingDate           string                 `json:"bookingDate"`
	Amount                Balance                `json:"amount"`
	Remitter              Account                `json:"remitter"`
	Deptor                Account                `json:"deptor"`
	Creditor              Creditor               `json:"creditor"`
	ValutaDate            string                 `json:"valutaDate"`
	DirectDebitCreditorID string                 `json:"directDebitCreditorId"`
	DirectDebitMandateID  string                 `json:"directDebitMandateId"`
	EndToEndReference     string                 `json:"endToEndReference"`
	NewTransaction        bool                   `json:"newTransaction"`
	RemittanceInfo        string                 `json:"remittanceInfo"`
	TransactionType       AccountTransactionType `json:"transactionType"`
}

type Creditor struct {
	HolderName string `json:"holderName"`
	IBAN       string `json:"iban"`
	BIC        string `json:"bic"`
}

type AccountTransactionType struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

type Depots struct {
	Paging Paging  `json:"paging"`
	Values []Depot `json:"values"`
}

type Depot struct {
	DepotID                    string   `json:"depotId"`
	DepotDisplayID             string   `json:"depotDisplayId"`
	ClientID                   string   `json:"clientId"`
	DepotType                  string   `json:"depotType"`
	DefaultSettlementAccountID string   `json:"defaultSettlementAccountId"`
	SettlementAccountIDs       []string `json:"settlementAccountIds"`
	TargetMarket               string   `json:"targetMarket"`
}

type DepotPositions struct {
	Paging              Paging                   `json:"paging"`
	AggregatedPositions AggregatedDepotPositions `json:"aggregated"`
	Values              []DepotPosition          `json:"values"`
}

type AggregatedDepotPositions struct {
	Depot                      Depot   `json:"depot"`
	PrevDayValue               Balance `json:"prevDayValue"`
	CurrentValue               Balance `json:"currentValue"`
	PurchaseValue              Balance `json:"purchaseValue"`
	ProfitLossPurchaseAbs      Balance `json:"profitLossPurchaseAbs"`
	ProfitLossPurchaseRel      string  `json:"profitLossPurchaseRel"`
	ProfitLossPrevDayAbs       Balance `json:"profitLossPrevDayAbs"`
	ProfitLossPrevDayRel       string  `json:"profitLossPrevDayRel"`
	ProfitLossPrevDayTotalAbs  Balance `json:"profitLossPrevDayTotalAbs"`
	PurchaseValuesAlterable    bool    `json:"purchaseValuesAlterable"`
	ProfitLossPrevDayTotalRel  string  `json:"profitLossPrevDayTotalRel"`
	ProfitLossPurchaseTotalAbs Balance `json:"profitLossPurchaseTotalAbs"`
	ProfitLossPurchaseTotalRel string  `json:"profitLossPurchaseTotalRel"`
}

type DepotPosition struct {
	DepotID                   string  `json:"depotId"`
	PositionID                string  `json:"positionId"`
	WKN                       string  `json:"wkn"`
	CustodyType               string  `json:"custodyType"`
	Quantity                  Balance `json:"quantity"`
	AvailableQuantity         Balance `json:"availableQuantity"`
	CurrentPrice              Price   `json:"currentPrice"`
	PurchasePrice             Balance `json:"purchasePrice"`
	PrevDayPrice              Price   `json:"prevDayPrice"`
	CurrentValue              Balance `json:"currentValue"`
	PurchaseValue             Balance `json:"purchaseValue"`
	ProfitLossPurchaseAbs     Balance `json:"profitLossPurchaseAbs"`
	ProfitLossPurchaseRel     string  `json:"profitLossPurchaseRel"`
	ProfitLossPrevDayAbs      Balance `json:"profitLossPrevDayAbs"`
	ProfitLossPrevDayRel      string  `json:"profitLossPrevDayRel"`
	ProfitLossPrevDayTotalAbs Balance `json:"profitLossPrevDayTotalAbs"`
	Version                   string  `json:"version"`
	Hedgeability              string  `json:"hedgeability"`
	AvailableQuantityToHedge  Balance `json:"availableQuantityToHedge"`
	CurrentPriceDeterminable  bool    `json:"currentPriceDeterminable"`
	HasIntraDayExecutedOrder  bool    `json:"hasIntraDayExecutedOrder"`
}

type Price struct {
	Price         Balance `json:"price"`
	PriceDateTime string  `json:"priceDateTime"`
	Venue         Venue   `json:"venue"`
}

type Venue struct {
	Name    string `json:"name"`
	VenueID string `json:"venueId"`
	Country string `json:"country"`
	Type    string `json:"type"`
}

type DepotTransactions struct {
	Paging Paging             `json:"paging"`
	Values []DepotTransaction `json:"values"`
}

type DepotTransaction struct {
	TransactionID        string     `json:"transactionId"`
	BookingStatus        string     `json:"bookingStatus"`
	BookingDate          string     `json:"bookingDate"`
	BusinessDate         string     `json:"businessDate"`
	Quantity             Balance    `json:"quantity"`
	InstrumentID         string     `json:"instrumentId"`
	Instrument           Instrument `json:"instrument"`
	ExecutionPrice       Price      `json:"executionPrice"`
	TransactionValue     Balance    `json:"transactionValue"`
	TransactionDirection string     `json:"transactionDirection"`
	TransactionType      string     `json:"transactionType"`
}

type Instrument struct {
	InstrumentID string               `json:"instrumentId"`
	WKN          string               `json:"wkn"`
	ISIN         string               `json:"isin"`
	Mnemonic     string               `json:"mnemonic"`
	Name         string               `json:"name"`
	ShortHand    string               `json:"shortName"`
	StaticData   StaticInstrumentData `json:"staticData"`
}

type StaticInstrumentData struct {
	Notation               string `json:"notation"`
	Currency               string `json:"currency"`
	InstrumentType         string `json:"instrumentType"`
	PriipsRelevant         bool   `json:"priipsRelevant"`
	KIDAvailable           bool   `json:"kidAvailable"`
	ShippingWaiverRequired bool   `json:"shippingWaiverRequired"`
	FundRedemptionLimited  bool   `json:"fundRedemptionLimited"`
	SavingsPlanEligibility string `json:"savingsPlanEligibility"`
}
