package ftx

import (
	"strconv"
	"time"
)

type PositionsResponse struct {
	Success bool       `json:"success"`
	Result  []Position `json:"result"`
}

type Position struct {
	Side         string  `json:"side"`
	Future       string  `json:"future"`
	NotionalSize float64 `json:"cost"`
	PositionSize float64 `json:"size"`
	UPNL         float64 `json:"unrealizedPnl"`
	EntryPrice   float64 `json:"entryPrice"`
}

func (f *Client) GetPosition() ([]Position, error) {
	var positions PositionsResponse

	var out []Position

	resp, err := f.get("positions?showAvgPrice=true", []byte(""))
	if err != nil {
		return out, err
	}
	err = processResponse(resp, &positions)

	for _, v := range positions.Result {
		if v.NotionalSize > 0.1 || v.NotionalSize < -0.1 {
			out = append(out, v)
		}
	}

	return out, err
}

type FillResponse struct {
	Success bool   `json:"success"`
	Result  []Fill `json:"result"`
}

type Fill struct {
	Fee         float64   `json:"fee,omitempty"`
	FeeCurrency string    `json:"fee_currency,omitempty"`
	Future      string    `json:"future,omitempty"`
	Id          int       `json:"id,omitempty"`
	OrderId     int       `json:"order_id,omitempty"`
	TradeId     int       `json:"trade_id,omitempty"`
	Price       float64   `json:"price,omitempty"`
	Side        string    `json:"side,omitempty"`
	Size        float64   `json:"size,omitempty"`
	Time        time.Time `json:"time,omitempty"`
}

func (f *Client) GetFills(st, et int64) ([]Fill, error) {
	var fills FillResponse

	var out []Fill

	resp, err := f.get("fills?start_time="+strconv.FormatInt(st, 10)+
		"&end_time="+strconv.FormatInt(et, 10), []byte(""))
	if err != nil {
		return out, err
	}
	err = processResponse(resp, &fills)

	return fills.Result, err
}
