package ftx

import (
	"log"
	"strconv"
	"time"
)

type FundingPaymentsResponse struct {
	Success bool              `json:"success"`
	Result  []FundingPayments `json:"result"`
}

type FundingPayments struct {
	Future  string    `json:"future,omitempty"`
	Id      int       `json:"id,omitempty"`
	Payment float64   `json:"payment,omitempty"`
	Time    time.Time `json:"time,omitempty"`
}

func (f *Client) GetFundingPayments(ticker string, st, et int64) ([]FundingPayments, error) {
	var fpr FundingPaymentsResponse

	var s string = "?"
	if ticker != "" {
		s = "?future=" + ticker + "&"
	}

	resp, err := f.get(
		"funding_payments"+s+
			"start_time="+strconv.FormatInt(st, 10)+
			"&end_time="+strconv.FormatInt(et, 10),
		[]byte(""))
	if err != nil {
		log.Println("ERROR GetFundingPayments", err)
		return nil, err
	}

	err = processResponse(resp, &fpr)

	return fpr.Result, nil
}

//GET /funding_payments

type FundingRatesResponse struct {
	Success bool           `json:"success"`
	Result  []FundingRates `json:"result"`
}

type FundingRates struct {
	Future string    `json:"future,omitempty"`
	Rate   float64   `json:"rate,omitempty"`
	Time   time.Time `json:"time,omitempty"`
}

func (p *Client) GetFundingRates(ticker string, st, et int64) ([]FundingRates, error) {
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
