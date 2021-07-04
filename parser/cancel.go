package parser

import (
	"errors"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

/*
cancel all
cancel buy btc-perp
cancel sell btc-perp limit
*/

type Cancel struct {
	lo     bool
	so     bool
	Ticker []string
}

func ParseCancel(tl []Token) (*Cancel, error) {
	var cancle Cancel
	cancle.Ticker = make([]string, 0)

	if len(tl) == 2 && tl[1].Text == "all" {
		return &cancle, nil
	}

	for _, v := range tl[1:] {
		switch v.Type {
		case FLAG:
			switch v.Text {
			case "limit":
				cancle.lo = true
			case "stop":
				cancle.so = true
			default:
				return nil, errors.New(v.Text + " Is not attribute for cancelation, try -limit or -stop")
			}
		case VARIABLE:
			cancle.Ticker = append(cancle.Ticker, v.Text)
		default:
			return nil, errors.New(v.Type.String() + " Is not a legit input for cancel")
		}
	}

	return &cancle, nil
}

func (c *Cancel) Evaluate(f *ftx.Client) error {
	if len(c.Ticker) == 0 {
		_, err := f.DeleteOrders("", c.so, c.lo)
		return err
	}

	for _, v := range c.Ticker {
		_, err := f.DeleteOrders(v, c.so, c.lo)
		if err != nil {
			return err
		}
	}
	return nil
}
