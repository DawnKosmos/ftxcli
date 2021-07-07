package ftx

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type AccountResponse struct {
	Success bool    `json:"success"`
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

type GetWalletBalanceResponse struct {
	Success bool   `json:"success"`
	Result  []Coin `json:"result"`
}

type Coin struct {
	Coin     string  `json:"coin,omitempty"`
	Free     float64 `json:"free,omitempty"`
	Total    float64 `json:"total,omitempty"`
	UsdValue float64 `json:"usdValue,omitempty"`
}

func (f *Client) GetWalletBalance() ([]Coin, error) {
	var ff GetWalletBalanceResponse
	resp, err := f.get("wallet/balances", []byte(""))
	if err != nil {
		return nil, err
	}
	err = processResponse(resp, &ff)
	var out []Coin
	for _, v := range ff.Result {
		if v.UsdValue > 1 {
			out = append(out, v)
		}
	}
	return out, err
}

type DepositHistoryResponse struct {
	Success bool      `json:"success"`
	Result  []Deposit `json:"result"`
}

type Deposit struct {
	Coin   string    `json:"coin,omitempty"`
	Fee    float64   `json:"fee,omitempty"`
	Id     int       `json:"id,omitempty"`
	Size   float64   `json:"size,omitempty"`
	Status string    `json:"status,omitempty"`
	Time   time.Time `json:"time,omitempty"`
	Notes  string    `json:"notes,omitempty"`
}

type Withdraw Deposit

func (f *Client) GetDepositHistory(st, et int64) ([]Deposit, error) {
	var fr DepositHistoryResponse
	resp, err := f.get(
		"wallet/deposits?start_time="+strconv.FormatInt(st, 10)+
			"&end_time="+strconv.FormatInt(et, 10),
		[]byte(""))
	if err != nil {
		log.Println("ERROR OHCLV FTX", err)
		return nil, err
	}
	err = processResponse(resp, &fr)
	return fr.Result, err
}

func (f *Client) GetWithdrawHistory(st, et int64) ([]Withdraw, error) {
	var fr DepositHistoryResponse
	resp, err := f.get(
		"wallet/withdrawals?start_time="+strconv.FormatInt(st, 10)+
			"&end_time="+strconv.FormatInt(et, 10),
		[]byte(""))
	if err != nil {
		log.Println("ERROR OHCLV FTX", err)
		return nil, err
	}
	err = processResponse(resp, &fr)
	var out []Withdraw
	return out, err
}

type TransferPayload struct {
	Coin        string  `json:"coin,omitempty"`
	Size        float64 `json:"size,omitempty"`
	Source      string  `json:"source,omitempty"`
	Destination string  `json:"destination,omitempty"`
}

type TransferResponse struct {
	Success bool `json:"success"`
	Result  struct {
		Id     int       `json:"id,omitempty"`
		Coin   string    `json:"coin,omitempty"`
		Size   float64   `json:"size,omitempty"`
		Time   time.Time `json:"time,omitempty"`
		Status string    `json:"status,omitempty"`
	} `json:"result,omitempty"`
}

func (f *Client) TransferBetweenAccounts(Coin string, Size float64, src, destination string) (TransferResponse, error) {
	var newOrderResponse TransferResponse
	requestBody, err := json.Marshal(TransferPayload{
		Coin:        Coin,
		Size:        Size,
		Source:      src,
		Destination: destination})
	if err != nil {
		log.Printf("Error PlaceOrder %v", err)
		return newOrderResponse, err
	}
	resp, err := f.post("orders", requestBody)
	if err != nil {
		log.Printf("Error PlaceOrder %v", err)
		return newOrderResponse, err
	}
	err = processResponse(resp, &newOrderResponse)
	return newOrderResponse, err
}
