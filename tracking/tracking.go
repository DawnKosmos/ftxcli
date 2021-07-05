package parser

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

type Month struct {
	Accountname string
	Data        []Day
	Last        []Checked
}

type CheckType int

const (
	FILLS CheckType = iota
	BALANCE
	FUNDINGPAYMENTS
	WITHDRAWDEPOSIT
)

type Checked struct {
	Type CheckType `json:"type,omitempty"`
	Time time.Time `json:"time,omitempty"`
}

type Day struct {
	Day     int
	Fills   []ftx.Fill
	FP      []ftx.FundingPayments
	WD      []WithDrawsDeposits
	Balance ftx.Account
}

type WithDrawsDeposits struct {
}

type Balance struct {
}
