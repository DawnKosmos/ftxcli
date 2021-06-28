package parser

import (
	"errors"
	"strconv"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

type OrderType int

const (
	SIMPLE OrderType = iota
	CASCADE
)

type Order struct {
	Side    string
	Ticker  string
	A       Amount
	P       Price
	Balance float64
}

type AmountType int

const (
	COIN AmountType = iota
	FIAT
	ACCOUNTSIZE
)

type Amount struct {
	Type AmountType
	Val  float64
}

func ParseOrder(Side string, tl []Token) (*Order, error) {
	var order Order
	order.Side = Side
	var err error

	if tl[0].Type == VARIABLE {
		order.Ticker = tl[0].Text
	}

	var amount Amount

	switch tl[1].Type {
	case FLOAT:
		amount.Type = COIN
	case UFLOAT:
		amount.Type = FIAT
	case PERCENT:
		amount.Type = ACCOUNTSIZE
	default:
		return nil, errors.New(tl[1].Type.String() + "is not a correct amount with value" + tl[1].Text)
	}

	amount.Val, err = strconv.ParseFloat(tl[1].Text, 64)

	if err != nil {
		return nil, err
	}

	order.A = amount

	order.P, err = ParsePrice(tl[2:])

	if err != nil {
		return &order, err
	}

	return &order, nil
}

type PriceType int

const (
	PRICE PriceType = iota
	DIFFERENCE
	PERCENTPRICE
)

type PriceSource int

type Price struct {
	Type       PriceType
	PC         string
	Duration   int64
	IsLaddered [2]bool
	//0,0 -> no, 1,0 -> laddered; 1,1 -> exponential laddered
	Values [3]float64
}

func ParsePrice(tl []Token) (p Price, err error) {

	p.PC = "market"
	if tl[0].Type == SOURCE {
		p.PC = tl[0].Text

		switch tl[0].Text {
		case "high":
			p.PC = "high"
		case "low":
			p.PC = "low"
		default:
			return p, errors.New(tl[0].Type.String() + "This Source does not exist with value" + tl[0].Text)
		}

		if tl[1].Type != DURATION {
			return p, errors.New("After -low, -high you need to provide a duration")
		}

		ss := tl[1].Text
		n, err := strconv.Atoi(ss[:len(ss)-1])
		if err != nil {
			return p, err
		}

		switch ss[len(ss)-1] {
		case 'h':
			n *= 3600
		case 'm':
			n *= 60
		case 'd':
			n *= 3600 * 24
		default:
			return p, errors.New(tl[1].Text + "I dont know how you fucked that up")
		}
		p.Duration = int64(n)

		tl = tl[2:]
	}

	switch tl[0].Type {
	case FLOAT:
		p.Type = PRICE
	case DFLOAT:
		p.Type = DIFFERENCE
	case PERCENT:
		p.Type = PERCENTPRICE
	case FLAG:
		err = ParsePriceFlag(tl[1:], &p, tl[0].Text)
		return p, err
	default:
		return p, errors.New(tl[0].Type.String() + "is not a correct Price with value" + tl[0].Text)
	}

	p.Values[0], err = strconv.ParseFloat(tl[1].Text, 64)
	if err != nil {
		return p, err
	}

	return p, nil
}

func ParsePriceFlag(tl []Token, p *Price, flag string) (err error) {

	switch flag {
	case "l":
		p.IsLaddered = [2]bool{true, false}
	case "le":
		p.IsLaddered = [2]bool{true, true}
	default:
		return errors.New("This Flag is not supported: " + flag)
	}

	if len(tl) < 3 {
		return errors.New("Not enough Arguments for a laddered order")
	}

	if tl[1].Type == FLOAT {
		num, err := strconv.Atoi(tl[0].Text)
		if err != nil {
			return err
		}

		if num > 20 || num < 2 {
			return errors.New("The number of laddering orders is to big or small with" + tl[1].Text)
		}

		p.Values[0] = float64(num)
	}

	if tl[1].Type != tl[2].Type {
		return errors.New("Num 1 and num 2 musst be from same typ")
	}

	switch tl[1].Type {
	case FLOAT:
		p.Type = PRICE
	case DFLOAT:
		p.Type = DIFFERENCE
	case PERCENT:
		p.Type = PERCENTPRICE
	default:
		return errors.New(tl[1].Type.String() + " is not a legit input for a laddered order")
	}
	v1, err := strconv.ParseFloat(tl[1].Text, 64)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	v2, err := strconv.ParseFloat(tl[2].Text, 64)

	p.Values[1], p.Values[2] = v1, v2

	// THEORETICAL PARSE A STOPS and OTHER EXTRA STUFF HERE
	return nil
}

func (o *Order) Evaluate(f ftx.Client) error {
	size, err := o.A.Evaluate(f, o.Ticker)
	if err != nil {
		return err
	}

}

func (a *Amount) Evaluate(f ftx.Client, ticker string) (size float64, err error) {
	switch a.Type {
	case COIN:
		return a.Val, nil
	case FIAT:
		m, err := f.GetMarket(ticker)
		if err != nil {
			return a.Val, err
		}
		temp := (m.Ask + m.Bid + m.Last) / 3
		return a.Val / temp, nil
	case ACCOUNTSIZE:
		m, err := f.GetMarket(ticker)
		if err != nil {
			return a.Val, err
		}
		temp := (m.Ask + m.Bid + m.Last) / 3
		account, err := f.GetAccount()
		if err != nil {
			return a.Val, err
		}
		az := account.FreeCollateral * a.Val / 100

		return az / temp, nil
	default:
		return
	}
}

func (p *Price) Evaluate(f ftx.Client, side string, ticker string, size float64) (err error) {
	if p.Type == PRICE {
		_, err = f.SetOrder(ticker, side, p.Values[0], size, "limit", false)
		return err
	}

	var mp float64
	if p.PC == "market" {
		m, err := f.GetMarket(ticker)
		if err != nil {
			return err
		}
		mp = (m.Ask + m.Bid + m.Last) / 3
	} else {
		mp, err = f.GetPriceSource(p.PC, ticker, p.Duration)
	}

	switch p.Type {
	case DIFFERENCE:

	case PERCENTPRICE:
	}
}

func (p *Price) EvaluateDifference(f ftx.Client, side string, ticker string, size float64, mp float64) error {
	var factor float64
	if side == "sell" {
		factor = -1.0
	}

	if !p.IsLaddered[0] {
		_, err := f.SetOrder(ticker, side, mp-p.Values[0]*factor, size, "limit", false)
		return err
	}

	p1, p2 := mp-p.Values[1]*factor, mp-p.Values[2]*factor
	plo := GetPricesLadderedOrder(p.IsLaddered[1], p.Values[0], p1, p2)

	for _, v := range plo {
		_, err := f.SetOrder(ticker, side, v[0], v[1], "limit", false)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetPricesLadderedOrder(exponential bool, split, p1, p2 float64) [][2]float64 {
	k := (p1 - p2) / split

	sum := split * (split + 1) / 20
	var fn func(iterate int, parti float64) float64
	if exponential {
		fn = func(iterate int, parti float64) float64 {
			return (float64(iterate+1) / parti) / sum
		}
	} else {
		fn = func(iterate int, parti float64) float64 {
			return 1 / parti
		}
	}

	var o [][2]float64
	for i := 0; i < int(split); i++ {
		o = append(o, [2]float64{p1 + k*float64(i), fn(i, split)})
	}
	return o
}
