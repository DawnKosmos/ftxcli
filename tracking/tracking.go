package tracking

import (
	"time"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

/*
Tracking soll alles was auf dem Account passiert tracken und in einer Datei speichern.
Diese Datei kann später benutzt werden um Statistiken zu erstellen und .csv dateien zu generieren
bzw. Statistiken auswerten



Daten werden in JSON format getrackt
jeder tag ist ein eintrag. jeder monat ist ein file

Im .acc file muss eine genehmigung fürs tracken drin sein
dann wird ein file in einem order gestellt und gefragt für wie viel zeit zurückgeschaut werden soll

Ein webseite soll erstellt werden mit allen nützlichen informationen wenn gefordert wird

[]TAG
	Fills[]
	[]FundingPayments
	[]Withdraws/Deposits
	Balance
*/

type CheckType int

const (
	FILLS CheckType = iota
	BALANCE
	FUNDINGPAYMENTS
	WITHDRAWDEPOSIT
)

type Tracker struct {
	Username    string    `json:"username,omitempty"`
	Month       string    `json:"month,omitempty"`
	Date        []Day     `json:"date,omitempty"`
	LastChecked []Checked `json:"last_checked,omitempty"`
}

type Checked struct {
	Type CheckType `json:"type,omitempty"`
	Time time.Time `json:"time,omitempty"`
}

type Day struct {
	Day     int               `json:"day,omitempty"`
	Fills   Fills             `json:"fills,omitempty"`
	FP      FundingPayment    `json:"fp,omitempty"`
	WD      WithdrawsDeposits `json:"wd,omitempty"`
	Balance []AccountBalance  `json:"balance,omitempty"`
}

type FundingPayment struct {
	Data []ftx.FundingPayments `json:"data,omitempty"`
}

type AccountBalance struct {
	Data []AccountSnapshot `json:"data,omitempty"`
}

type AccountSnapshot struct {
	Total float64    `json:"total,omitempty"`
	Coins []ftx.Coin `json:"coins,omitempty"`
	Time  time.Time  `json:"time,omitempty"`
}

type WithdrawsDeposits struct {
	Deposits  []ftx.Deposit  `json:"deposits,omitempty"`
	Withdraws []ftx.Withdraw `json:"withdraws,omitempty"`
}

type Fills struct {
	Data []ftx.Fill `json:"data,omitempty"`
}
