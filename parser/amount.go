package parser

import (
	"github.com/DawnKosmos/ftxcmd/ftx"
)

type AmountType int

const (
	COIN AmountType = iota
	FIAT
	ACCOUNTSIZE
	POSITIONSIZE
)

type Amount struct {
	Type AmountType
	Val  float64
}

func (a *Amount) Evaluate(f *ftx.Client, ticker string) (size float64, err error) {

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
	case POSITIONSIZE:
		pz, err := f.GetPosition()
		if err != nil {
			return a.Val, err
		}
		for _, v := range pz {
			if v.Future == ticker {
				return v.PositionSize, nil
			}
		}
	default:
		return
	}
	return a.Val, nil
}
