package parser

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

type SumOptions int

const (
	NONE SumOptions = iota
	ASCENDING
	DESCENDING
)

type FundingRates struct {
	Ticker []string
	Time   int64

	Summarize bool
}

func ParseFundingRates(tl []Token) (*FundingRates, error) {
	fund := &FundingRates{
		Ticker:    []string{},
		Time:      36000,
		Summarize: false,
	}
	var err error

	for _, v := range tl {
		switch v.Type {
		case FLAG:
			switch v.Text {
			case "sum":
				fund.Summarize = true
			}
		case VARIABLE:
			fund.Ticker = append(fund.Ticker, v.Text)
		case DURATION:
			fund.Time, err = parseDuration(v.Text)
			if err != nil {
				return fund, err
			}
		case FLOAT:
			ff, err := strconv.ParseFloat(v.Text, 64)
			if err != nil {
				return nil, err
			}
			fund.Time = int64(ff) * 3600
		}
	}
	return fund, nil
}

func parseDuration(ss string) (int64, error) {
	n, err := strconv.Atoi(ss[:len(ss)-1])
	if err != nil {
		return 0, err
	}
	switch ss[len(ss)-1] {
	case 'h':
		n *= 3600
	case 'm':
		n *= 60
	case 'd':
		n *= 3600 * 24
	default:
		return 0, errors.New(ss + " I dont know how you fucked that up")
	}

	return int64(n), nil
}

func (f *FundingRates) Evaluate(c *ftx.Client) (err error) {
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

func PrintFunding(summarize bool, fp ...[]ftx.FundingRates) {
	printfr := make([][]ftx.FundingRates, len(fp[0]))
	for i, v := range fp[0] {
		printfr[i] = []ftx.FundingRates{v}
	}
	if len(fp) == 1 {

	} else {
		for _, v := range fp[1:] {
			for i, vv := range v {
				printfr[i] = append(printfr[i], vv)
			}
		}
	}
	fmt.Print("Ticker: \t")

	for _, v := range printfr[0] {
		fmt.Print(v.Future, " ")
	}

	if summarize {
		ff := make([]float64, len(printfr[0]), len(printfr[0]))

		for _, v := range printfr {
			for i, vv := range v {
				ff[i] = ff[i] + vv.Rate
			}
		}

		fmt.Print("\nSummarized\t")
		for _, v := range ff {
			ss := fmt.Sprintf("%.4f", float64(v*100))
			fmt.Print(ss, "\t")
		}
		return
	}

	fmt.Print("\n")
	for _, v := range printfr {
		fmt.Print(v[0].Time.Format("02.07.06 15"), "\t")
		for _, vv := range v {
			ff := fmt.Sprintf("%.4f", float64(vv.Rate*100))
			fmt.Print(ff, "\t")
		}
		fmt.Print("\n")
	}

	return
}
