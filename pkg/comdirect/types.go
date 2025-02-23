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
