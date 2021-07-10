package parser

import (
	"errors"
	"strconv"
	"time"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

/*

funding -position //funding rate der aktuellen positionen
funding -highest 20 //funding der highest 20 coins
funding 10
funding 6h
funding -sum 1d
funding btc eth 10
*/

type FundingType int

const (
	PAYMENTS FundingType = iota
	POSITIONS
)

type Funding struct {
	ft        FundingType
	Ticker    []string
	Time      int64
	Summarize bool
}

func ParseFunding(tl []Token) (*Funding, error) {
	var fund Funding
	fund.ft = PAYMENTS
	var err error

	fund.Time = 36000
	for _, v := range tl {
		switch v.Type {
		case FLAG:
			switch v.Text {
			case "sum":
				fund.Summarize = true
			case "position":
				fund.ft = POSITIONS
			}
		case VARIABLE:
			fund.Ticker = append(fund.Ticker, v.Text)
		case DURATION:
			fund.Time, err = parseDuration(v.Text)
			if err != nil {
				return &fund, err
			}
		case FLOAT:
			ff, err := strconv.ParseFloat(v.Text, 64)
			if err != nil {
				return nil, err
			}
			fund.Time = int64(ff) * 3600
		default:
			return nil, errors.New(v.Type.String() + " Type not supported in Funding")
		}
	}
	return &fund, nil
}

/*
	Payments: Show Payed fees for a time period + have option to sum + sort by pair
	Positions: Showing the funding fees since the position got created
	General: sowing funding of coins, sum up
*/

func (f *Funding) Evaluate(c *ftx.Client) (err error) {
	switch f.ft {
	case PAYMENTS:
		err = EvaluatePayments(f, c)
	case POSITIONS:
		err = EvaluatePositions(f, c)
	}

	return err
}

func EvaluatePayments(f *Funding, c *ftx.Client) error {
	now := time.Now().Unix()
	if len(f.Ticker) == 0 {
		fp, err := c.GetFundingPayments("", now-f.Time, now)
		if err != nil {
			return err
		}
		err = PrintFundingPayments(f.Summarize, fp)
		return err
	}

	var fpr [][]ftx.FundingPayments
	for _, v := range f.Ticker {
		fp, err := c.GetFundingPayments(v, now-f.Time, now)
		if err != nil {
			return err
		}

		fpr = append(fpr, fp)
	}
	err := PrintFundingPayments(f.Summarize, fpr...)
	return err

}

func EvaluatePositions(f *Funding, c *ftx.Client) error {
	p, err := c.GetPosition()
	if err != nil {
		return err
	}
	if len(p) == 0 {
		return errors.New("No existing Positions found")
	}

	for _, v := range p {
		if len(v.Future) <= 5 {
			continue
		}

		if v.Future[len(v.Future)-4:] == "perp" {
			f.Ticker = append(f.Ticker, v.Future)
		}
	}

	err = EvaluatePayments(f, c)

	return err
}

// PRINT THE FUNCTIONS
func PrintFundingPayments(summarize bool, fp ...[]ftx.FundingPayments) error {
	//printfr := make([][]ftx.FundingPayments, 0)
	TimePosition := make(map[time.Time][]ftx.FundingPayments)

	for _, v := range fp {
		for _, vv := range v {
			r, ok := TimePosition[vv.Time]
			if !ok {
				TimePosition[vv.Time] = []ftx.FundingPayments{vv}
			} else {
				TimePosition[vv.Time] = append(r, vv)
			}
		}
	}

	/*
		Perp list [btc-perp, eth-perp, xrp-perp]
		inputs sortiert werden reingeschickt und mit perp list verglichen
		output is ftx.Fundingpayments array wenn nicht vorhanden ist der wert nil
	*/

	return nil
}

type FundingPaymentArray []ftx.FundingPayments

func (a FundingPaymentArray) Len() int      { return len(a) }
func (a FundingPaymentArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FundingPaymentArray) Less(i, j int) bool {
	return a[i].Future[0] < a[j].Future[0]
}
