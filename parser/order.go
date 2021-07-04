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
	Side   string
	Ticker string
	A      Amount
	P      Price
}

func ParseOrder(Side string, tl []Token) (*Order, error) {
	var order Order
	order.Side = Side
	var amount Amount
	var err error

	if len(tl) == 0 {
		return nil, errors.New("Parse Order has to have an input and cant be empty")
	}

	if tl[0].Type == VARIABLE {
		order.Ticker = tl[0].Text
	}

	//The amount getting created
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

	//Price getting creating
	order.P, err = ParsePrice(tl[2:])

	if err != nil {
		return &order, err
	}

	return &order, nil
}

//Evaluate the ordersize and prices
func (o *Order) Evaluate(f *ftx.Client) error {
	//We first get the amount evaluate, it depends on the type of amount and price source, low, high or market
	size, err := o.A.Evaluate(f, o.Ticker)
	if err != nil {
		return err
	}

	//the price get evaluation also places the orders
	err = o.P.Evaluate(f, o.Side, o.Ticker, size)
	return err
}
