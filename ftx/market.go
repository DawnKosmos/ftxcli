package ftx

import "log"

/*
    "name": "BTC-0628",
     "baseCurrency": null,
     "quoteCurrency": null,
     "type": "future",
     "underlying": "BTC",
     "enabled": true,
     "ask": 3949.25,
     "bid": 3949,
     "last": 10579.52,
     "postOnly": false,
     "priceIncrement": 0.25,
     "sizeIncrement": 0.0001,
     "restricted": false
   }
*/

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

func (c *Client) GetMarket(ticker string) (MarketResponse, error) {
	var r MarketResponse

	resp, err := c.get("markets/"+ticker, []byte(""))
	if err != nil {
		log.Println("ERROR GetMarket", err)
		return r, err
	}

	err = processResponse(resp, &r)

	return r, nil
}

/*
func (p *Public) GetFundingRates(ticker string, st, et int64) ([]FundingRates, error) {
	var fr FundingRatesResponse

	resp, err := p.get(
		"funding_rates?future="+ticker+
			"&start_time="+strconv.FormatInt(st, 10)+
			"&end_time="+strconv.FormatInt(et, 10),
		[]byte(""))

	if err != nil {
		log.Println("ERROR OHCLV FTX", err)
		return nil, err
	}

	err = processResponse(resp, &fr)

	return fr.Result, nil

}
*/
