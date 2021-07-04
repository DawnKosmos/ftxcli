package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

//stop buy btc-perp [position, u100,0.1] -low 5h

type TriggerOrder struct {
	Side   string
	Ticker string
	Amount Amount

	PType      PriceType
	PC         string
	Duration   int64
	Value      float64
	ReduceOnly bool
}

func ParseStop(side string, tl []Token) (*TriggerOrder, error) {
	if len(tl) == 0 {

	}
	var tOrder TriggerOrder
	tOrder.Side = side
	var amount Amount

	var err error

	if tl[0].Type == VARIABLE {
		tOrder.Ticker = tl[0].Text
	}

	switch tl[1].Type {
	case FLOAT:
		amount.Type = COIN
	case UFLOAT:
		amount.Type = FIAT
	case PERCENT:
		amount.Type = ACCOUNTSIZE
	case POSITION:
		amount.Type = POSITIONSIZE
	default:
		return nil, errors.New(tl[1].Type.String() + "is not a correct amount with value" + tl[1].Text)
	}

	if amount.Type != ACCOUNTSIZE {
		amount.Val, err = strconv.ParseFloat(tl[1].Text, 64)

		if err != nil {
			return nil, err
		}
	}

	tOrder.Amount = amount

	tOrder.PC = "market"
	tl = tl[2:]

	if tl[0].Type == SOURCE {
		switch tl[0].Text {
		case "high":
			tOrder.PC = "high"
		case "low":
			tOrder.PC = "low"
		default:
			return &tOrder, errors.New(tl[0].Type.String() + "This Source does not exist with value" + tl[0].Text)
		}

		if tl[1].Type != DURATION {
			return &tOrder, errors.New("After -low, -high you need to provide a duration")
		}

		ss := tl[1].Text
		n, err := strconv.Atoi(ss[:len(ss)-1])
		if err != nil {
			return &tOrder, err
		}

		switch ss[len(ss)-1] {
		case 'h':
			n *= 3600
		case 'm':
			n *= 60
		case 'd':
			n *= 3600 * 24
		default:
			return &tOrder, errors.New(tl[1].Text + "I dont know how you fucked that up")
		}
		tOrder.Duration = int64(n)

		tl = tl[2:]
	}

	switch tl[0].Type {
	case FLOAT:
		tOrder.PType = PRICE
	case DFLOAT:
		tOrder.PType = DIFFERENCE
	case PERCENT:
		tOrder.PType = PERCENTPRICE
	default:
		return &tOrder, errors.New(tl[0].Type.String() + "is not a correct Price with value" + tl[0].Text)
	}

	tOrder.Value, err = strconv.ParseFloat(tl[1].Text, 64)
	if err != nil {
		return &tOrder, err
	}

	tl = tl[1:]

	if len(tl) != 0 {
		if tl[0].Type == FLAG && tl[0].Text == "ro" {
			tOrder.ReduceOnly = true
		}
	}

	return &tOrder, nil
}

func (t *TriggerOrder) Evaluate(f *ftx.Client) error {
	size, err := t.Amount.Evaluate(f, t.Ticker)
	if err != nil {
		return err
	}

	factor := -1.0
	if t.Side == "sell" {
		factor = 1.0
	}

	var mp float64

	if t.PC == "market" {
		m, err := f.GetMarket(t.Ticker)
		if err != nil {
			return err
		}
		mp = (m.Ask + m.Bid + m.Last) / 3
	} else {
		mp, err = f.GetPriceSource(t.PC, t.Ticker, t.Duration)
	}

	switch t.PType {
	case PRICE:
		mp = t.Value
	case DIFFERENCE:
		mp = mp - t.Value*factor
	case PERCENTPRICE:
		mp = mp - mp*t.Value/100*factor
	}

	p, err := f.SetTriggerOrder(t.Ticker, t.Side, mp, size, "stop", t.ReduceOnly)
	if err != nil {
		return err
	}
	fmt.Println(p.Result.Side, p.Result.Size, p.Result.OrderPrice)

	return nil
}
