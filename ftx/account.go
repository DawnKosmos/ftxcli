package ftx

type AccountResponse struct {
	Success bool
	Result  Account `json:"result"`
}

type Account struct {
	Username          string     `json:"username"`
	Collateral        float64    `json:"collateral"`
	FreeCollateral    float64    `json:"freeCollateral"`
	Leverage          float64    `json:"leverage"`
	MarginFraction    float64    `json:"marginFraction"`
	TotalAccountValue float64    `json:"totalAccountValue"`
	TotalPositionSize float64    `json:"totalPositionSize"`
	Positions         []Position `json:"positions"`
}

func (pr *Client) GetAccount() (Account, error) {
	var ff AccountResponse
	var out Account
	resp, err := pr.get("account", []byte(""))
	if err != nil {
		return out, err
	}
	err = processResponse(resp, &ff)
	if err != nil {
		return out, err
	}

	out = ff.Result
	var p []Position

	for _, v := range out.Positions {
		if v.NotionalSize > 0.1 || v.NotionalSize < -0.1 {
			p = append(p, v)
		}
	}

	out.Positions = p

	return out, nil
}
