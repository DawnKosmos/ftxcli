package tracking

import "github.com/DawnKosmos/ftxcmd/ftx"

type FundingPayment struct {
	Data []ftx.FundingPayments `json:"data,omitempty"`
}

func (f *FundingPayment) Add(v ftx.FundingPayments) {
	f.Data = append(f.Data, v)
}

func (d *Day) SummarizeFundingFees() {
	TickerFund := make(map[string]ftx.FundingPayments)

	for _, v := range d.FP.Data {
		vv, ok := TickerFund[v.Future]
		if !ok {
			TickerFund[v.Future] = vv
		} else {
			vv.Payment += v.Payment
			TickerFund[v.Future] = vv
		}
	}

	var fp []ftx.FundingPayments
	for _, v := range TickerFund {
		fp = append(fp, v)
	}
}
