package ftx

import (
	"log"
	"strconv"
	"time"
)

type MarketResponse struct {
	Success bool   `json:"success"`
	Result  Market `json:"result"`
}

type Market struct {
	Name           string  `json:"name"`
	Ask            float64 `json:"ask"`
	Bid            float64 `json:"bid"`
	Last           float64 `json:"last"`
	PriceIncrement float64 `json:"priceIncrement"`
	SizeIncrement  float64 `json:"sizeIncrement"`
}

func (c *Client) GetMarket(ticker string) (Market, error) {
	var r MarketResponse

	resp, err := c.get("markets/"+ticker, []byte(""))
	if err != nil {
		log.Println("ERROR GetMarket", err)
		return r.Result, err
	}

	err = processResponse(resp, &r)

	return r.Result, nil
}

type Chart struct {
	Success bool     `json:"success"`
	Result  []Candle `json:"result"`
}

type Candle struct {
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Open      float64   `json:"open"`
	Volume    float64   `json:"volume"`
	StartTime time.Time `json:"startTime"`
}

func (p *Client) GetHistoricalPrices(ticker string, res, st, et int64) ([]Candle, error) {
	var hp Chart
	resp, err := p.get(
		"markets/"+ticker+
			"/candles?resolution="+strconv.FormatInt(res, 10)+
			"&start_time="+strconv.FormatInt(st, 10)+
			"&end_time="+strconv.FormatInt(et, 10),
		[]byte(""))
	if err != nil {
		log.Println("ERROR OHCLV FTX", err)
		return nil, err
	}

	err = processResponse(resp, &hp)

	return hp.Result, nil
}

func (p *Client) GetPriceSource(src string, ticker string, duration int64) (out float64, err error) {
	if src == "market" {
		m, err := p.GetMarket(ticker)
		if err != nil {
			return 0, err
		}
		return (m.Ask + m.Bid + m.Last) / 3, nil
	}

	et := time.Now().Unix()
	st := et - duration
	var ch []Candle
	if duration < 3600 {
		ch, err = p.GetHistoricalPrices(ticker, 60, st, et)
		if err != nil {
			return
		}
	} else {
		ch, err = p.GetHistoricalPrices(ticker, 3600, st, et)
		if err != nil {
			return
		}
	}

	out = LowestHighestCandles(src, ch)

	return out, nil
}

func LowestHighestCandles(src string, ch []Candle) (out float64) {
	if src == "low" {
		out = ch[0].Low
		for _, c := range ch {
			if c.Low < out {
				out = c.Low
			}
		}
	} else {
		out = ch[0].High
		for _, c := range ch {
			if c.High > out {
				out = c.High
			}
		}
	}
	return out
}
