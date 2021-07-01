package parser

import (
	"errors"
	"strconv"
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
