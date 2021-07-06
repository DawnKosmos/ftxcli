package tracking

import (
	"errors"
	"fmt"
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
	Day     int `json:"day,omitempty"`
	Month   time.Month
	Year    int
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

func (f Fills) Add(v ftx.Fill) {
	f.Data = append(f.Data, v)
}

func (t *Tracker) Fill(f *ftx.Client, stt, ett time.Time) error {
	if stt.Month() != stt.Month() {
		return errors.New("Start Month and End Month aren't the same")
	}
	wb, err := f.GetWalletBalance()
	if err != nil {
		return err
	}

	st, et := stt.Unix(), ett.Unix()

	var total float64
	for _, v := range wb {
		total += v.UsdValue
	}
	fmt.Println(total)

	fills, err := f.GetFills(st, et)
	if err != nil {
		return err
	}

	for _, v := range fills {
		i := v.Time.Day()
		t.Date[i-1].Fills.Add(v)
	}

	return nil
}
