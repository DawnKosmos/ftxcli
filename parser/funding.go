package parser

import (
	"errors"
	"fmt"
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

func (f *Funding) Evaluate(c *ftx.Client, ws *WsAccount) (err error) {
	switch f.ft {
	case PAYMENTS:
		err = EvaluatePayments(f, c, ws)
	case POSITIONS:
		err = EvaluatePositions(f, c, ws)
	}

	return err
}

func EvaluatePayments(f *Funding, c *ftx.Client, ws *WsAccount) error {
	now := time.Now().Unix()
	if len(f.Ticker) == 0 {
		fp, err := c.GetFundingPayments("", now-f.Time, now)
		if err != nil {
			return err
		}
		if f.Summarize {
			if ws == nil {
				fmt.Print("Funding rate payed: ")
			} else {
				ws.AddToBuffer("Funding rate payed:\t")
			}
			var sum float64
			for _, v := range fp {
				sum += v.Payment
			}
			if ws == nil {
				fmt.Println("\t", sum)
			} else {
				ws.Write(fmt.Sprintf("\t %f", sum))
			}
			return nil
		}

		for _, vv := range fp {
			if ws == nil {
				fmt.Println(vv.Future, vv.Time.Format("02.07.06 15"), vv.Payment)
			} else {
				ws.Write(fmt.Sprintf("%s %s %f", vv.Future, vv.Time.Format("02.07.06 15"), vv.Payment))
			}
		}

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
	PrintFundingPayments(f.Summarize, ws, fpr...)
	return nil

}

func EvaluatePositions(f *Funding, c *ftx.Client, ws *WsAccount) error {
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

	err = EvaluatePayments(f, c, ws)

	return err
}

// PRINT THE FUNCTIONS
func PrintFundingPayments(summarize bool, ws *WsAccount, fp ...[]ftx.FundingPayments) {
	var msg []byte
	if summarize {
		ff := make([]float64, len(fp))

		if ws == nil {
			fmt.Print("Ticker: \t")
			for _, v := range fp {
				fmt.Print(v[0].Future, " ")
			}
			fmt.Print("\nSummarized\t")
		} else {
			msg = append(msg, []byte("Ticker: \t")...)
			for _, v := range fp {
				msg = append(msg, []byte(v[0].Future+" ")...)
			}
			msg = append(msg, []byte("\nSummarized\t")...)
		}

		for i, v := range fp {
			for _, vv := range v {
				ff[i] = ff[i] + vv.Payment
			}
		}

		for _, v := range ff {
			ss := fmt.Sprintf("%.4f", float64(v))
			fmt.Print(ss, "\t")
		}
		return
	}

	for _, v := range fp {
		fmt.Println("Ticker: \t", v[0].Future)
		for _, vv := range v {
			fmt.Println(vv.Time.Format("02.07.06 15"), vv.Payment)
		}
	}

}
