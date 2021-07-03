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
	GENERAL
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
			ss := v.Text
			n, err := strconv.Atoi(ss[:len(ss)-1])
			if err != nil {
				return nil, err
			}
			switch ss[len(ss)-1] {
			case 'h':
				n *= 3600
			case 'm':
				n *= 60
			case 'd':
				n *= 3600 * 24
			default:
				return &fund, errors.New(v.Text + "I dont know how you fucked that up")
			}
			fund.Time = int64(n)
		case FLOAT:
			ff, err := strconv.ParseFloat(v.Text, 64)
			if err != nil {
				return nil, err
			}
			fund.Time = int64(ff) * 3600
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
	case GENERAL:
		err = EvaluateGeneral(f, c)
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
		if v.Future[:len(v.Future)-4] == "perp" {
			f.Ticker = append(f.Ticker, v.Future)
		}
	}

	err = EvaluatePayments(f, c)

	return err
}

func EvaluateGeneral(f *Funding, c *ftx.Client) error {
	t := time.Now().Unix()
	if len(f.Ticker) == 0 {
		fp, err := c.GetFundingRates("", t, t)
		if err != nil {
			return err
		}
		PrintFunding(f.Summarize, fp)
		return nil
	}

	var fpr [][]ftx.FundingRates

	for _, v := range f.Ticker {
		fp, err := c.GetFundingRates(v, t-f.Time, t)
		if err != nil {
			return err
		}
		fpr = append(fpr, fp)

	}

	PrintFunding(f.Summarize, fpr...)

	return nil
}

// PRINT THE FUNCTIONS
func PrintFundingPayments(summarize bool, fp ...[]ftx.FundingPayments) error {
	for _, v := range fp {
		fmt.Println(v)
	}
	return nil
}

/*	mapfr := make(map[int64][]ftx.FundingRates)
	for _, v := range fp {
		for _, vv := range v {
			fpr, ok := mapfr[vv.Time.Unix()]
			if !ok {
				mapfr[vv.Time.Unix()] = []ftx.FundingRates{vv}
			} else {
				fpr = append(fpr, vv)
				mapfr[vv.Time.Unix()] = fpr
			}
		}
	}

	for k, v := range mapfr {
		fmt.Print(k, " ")
		for _, vv := range v {
			fmt.Print(vv.Future, vv.Rate, " ")
		}
		fmt.Print("\n")
	}
	return nil*/
