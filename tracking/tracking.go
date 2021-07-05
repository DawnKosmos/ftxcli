package parser

import "github.com/DawnKosmos/ftxcmd/ftx"

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

type Day struct {
	Day     int
	Fills   []Fill
	FP      []ftx.FundingPayments
	WD      []WithDrawsDeposits
	Balance ftx.Account
}

type Fill struct {
}

type WithDrawsDeposits struct {
}

type Balance struct {
}
