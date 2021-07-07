package tracking

import "github.com/DawnKosmos/ftxcmd/ftx"

type Fills struct {
	Data []ftx.Fill `json:"data,omitempty"`
}

func (f *Fills) Add(v ftx.Fill) {
	f.Data = append(f.Data, v)
}

func (f *Day) GetDailyPaidFees() []ftx.Fill {
	TickerFill := make(map[string]ftx.Fill)

	for _, v := range f.Fills.Data {
		vv, ok := TickerFill[v.Future]
		if !ok {
			TickerFill[v.Future] = vv
		} else {
			vv.Fee += v.Fee
			vv.Size += v.Size
			TickerFill[v.Future] = vv
		}
	}
	var fp []ftx.Fill
	for _, v := range TickerFill {
		fp = append(fp, v)
	}
	return fp
}
