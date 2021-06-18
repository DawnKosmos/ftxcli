package ftx

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
