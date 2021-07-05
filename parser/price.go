package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

type PriceSource int

type PriceType int

const (
	PRICE PriceType = iota
	DIFFERENCE
	PERCENTPRICE
)

type Price struct {
	Type       PriceType
	PC         string
	Duration   int64
	IsLaddered [2]bool
	//0,0 -> no, 1,0 -> laddered; 1,1 -> exponential laddered
	Values [3]float64
}

func ParsePrice(tl []Token) (p Price, err error) {
	if len(tl) == 0 {
		return p, errors.New("An order needs arguments")
	}

	p.PC = "market"
	if tl[0].Type == SOURCE {

		switch tl[0].Text {
		case "high":
			p.PC = "high"
		case "low":
			p.PC = "low"
		default:
			return p, errors.New(tl[0].Type.String() + "This Source does not exist with value" + tl[0].Text)
		}

		if len(tl) < 2 || tl[1].Type != DURATION {
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

	p.Values[0], err = strconv.ParseFloat(tl[0].Text, 64)
	if err != nil {
		return p, err
	}

	return p, nil
}

func ParsePriceFlag(tl []Token, p *Price, flag string) (err error) {
	if len(tl) > 3 {
		return errors.New("Not enough arguments")
	}
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

	if tl[0].Type == FLOAT {
		num, err := strconv.Atoi(tl[0].Text)
		if err != nil {
			return err
		}

		if num > 20 || num < 2 {
			return errors.New("The number of laddering orders is to big or small with" + tl[1].Text)
		}

		p.Values[0] = float64(num)
	} else {
		return errors.New("PARSE ERROR laddered TYPE should be FLOAT is " + tl[1].Type.String())
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

func (p *Price) Evaluate(f *ftx.Client, side string, ticker string, size float64) (err error) {
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
	case PRICE:
		err = p.EvaluatePrice(f, side, ticker, size, mp)
	case DIFFERENCE:
		err = p.EvaluateDifference(f, side, ticker, size, mp)
	case PERCENTPRICE:
		err = p.EvaluatePercentual(f, side, ticker, size, mp)
	}

	return err
}

func (p *Price) EvaluatePrice(f *ftx.Client, side string, ticker string, size float64, mp float64) error {

	if !p.IsLaddered[0] {
		_, err := f.SetOrder(ticker, side, p.Values[0], size, "limit", false)
		return err
	}

	p1, p2 := p.Values[1], p.Values[2]
	plo := GetPricesLadderedOrder(p.IsLaddered[1], p.Values[0], p1, p2)

	for _, v := range plo {
		_, err := f.SetOrder(ticker, side, v[0], size*v[1], "limit", false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Price) EvaluateDifference(f *ftx.Client, side string, ticker string, size float64, mp float64) error {
	factor := 1.0
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
		_, err := f.SetOrder(ticker, side, v[0], size*v[1], "limit", false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Price) EvaluatePercentual(f *ftx.Client, side string, ticker string, size float64, mp float64) error {
	factor := 1.0
	if side == "sell" {
		factor = -1.0
	}
	//var err error

	if !p.IsLaddered[0] {
		_, err := f.SetOrder(ticker, side, mp-mp*p.Values[0]/100*factor, size, "limit", false)
		return err
	}

	p11, p22 := mp*p.Values[1]/100, mp*p.Values[2]/100
	p1, p2 := mp-p11*factor, mp-p22*factor
	plo := GetPricesLadderedOrder(p.IsLaddered[1], p.Values[0], p1, p2)

	for _, v := range plo {
		p, err := f.SetOrder(ticker, side, v[0], size*v[1], "limit", false)
		if err != nil {
			return err
		}
		fmt.Println(p.Result.Side, p.Result.Size, p.Result.Price)

	}

	return nil
}

func GetPricesLadderedOrder(exponential bool, split, p1, p2 float64) [][2]float64 {
	b := (p2 - p1) / split
	k := b * split / (split - 1)
	//k := (p2 - p1 + b) / split

	sum := (split + 1) / 2
	var fn func(iterate int) float64

	// i(i)
	if exponential {
		fn = func(iterate int) float64 {
			return (float64(iterate+1) / split) / sum
		}
	} else {
		fn = func(iterate int) float64 {
			return 1 / split
		}
	}

	var o [][2]float64
	for i := 0; i < int(split); i++ {
		o = append(o, [2]float64{p1 + k*float64(i), fn(i)})
	}
	return o
}

func harmonicSum(n int) float64 {
	var sum float64
	for i := 0; i < n; i++ {
		sum += 1 / (float64(i) + 1)
	}
	return sum
}
