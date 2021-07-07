package tracking

import (
	"errors"
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
	ALL CheckType = iota
	FILLS
	BALANCE
	FUNDINGPAYMENTS
	WITHDRAWDEPOSIT
)

type Tracker struct {
	Username    string     `json:"username,omitempty"`
	Month       string     `json:"month,omitempty"`
	Year        int        `json:"year,omitempty"`
	Date        []Day      `json:"date,omitempty"`
	LastChecked [5]Checked `json:"last_checked,omitempty"`
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
	Balance AccountBalance    `json:"balance,omitempty"`
}

func (t *Tracker) Fill(f *ftx.Client, stt, ett time.Time) error {
	st, et := stt.Unix(), ett.Unix()
	todayDay := time.Now().Day()
	if stt.Month() != stt.Month() {
		return errors.New("Start Month and End Month aren't the same")
	}

	if f.Subaccount == "" {
		t.Username = "main"
	} else {
		t.Username = f.Subaccount
	}

	//Get Total Wallet Balance
	wb, err := f.GetWalletBalance()
	if err != nil {
		return err
	}
	t.Date[todayDay-1].Balance.Add(wb)
	//Get Fills
	fills, err := f.GetFills(st, et)
	if err != nil {
		return err
	}

	for _, v := range fills {
		i := v.Time.Day()
		t.Date[i-1].Fills.Add(v)
	}
	//withdraw deposits
	withdraws, err := f.GetWithdrawHistory(st, et)
	if err != nil {
		return err
	}
	deposits, err := f.GetDepositHistory(st, et)
	if err != nil {
		return err
	}
	for _, v := range withdraws {
		i := v.Time.Day()
		t.Date[i-1].WD.AddWithdraw(v)
	}
	for _, v := range deposits {
		i := v.Time.Day()
		t.Date[i-1].WD.AddDeposit(v)
	}
	//Fundingpayments
	fp, err := f.GetFundingPayments("", st, et)
	for _, v := range fp {
		i := v.Time.Day()
		t.Date[i-1].FP.Add(v)
	}

	t.LastChecked[ALL] = Checked{ALL, time.Now()}
	return nil
}
