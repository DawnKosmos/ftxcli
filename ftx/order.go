package ftx

import (
	"encoding/json"
	"log"
	"time"
)

type OpenOrdersRequest struct {
	Success bool    `json:"success"`
	Result  []Order `json:"result"`
}

type Order struct {
	ID         int64     `json:"id"`
	Ticker     string    `json:"market"`
	Side       string    `json:"side"`
	Size       float64   `json:"size"`
	Price      float64   `json:"price"`
	ReduceOnly bool      `json:"reduceOnly"`
	FilledSize float64   `json:"filledSize"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (f *Client) GetOpenOrders(ticker ...string) ([]Order, error) {
	var oor OpenOrdersRequest

	s := "orders"
	if len(ticker) != 0 {
		s += "?market=" + ticker[0]
	}

	resp, err := f.get(s, []byte(""))
	if err != nil {
		return nil, err
	}
	err = processResponse(resp, &oor)

	out := oor.Result

	return out, err
}

/*
{
  "market": "XRP-PERP",
  "side": "sell",
  "price": 0.306525,
  "type": "limit",
  "size": 31431.0,
  "reduceOnly": false,
  "ioc": false,
  "postOnly": false,
  "clientId": null
}
*/

type NewOrderResponse struct {
	Success bool  `json:"success"`
	Result  Order `json:"result"`
}

type NewOrder struct {
	Market                  string  `json:"market"`
	Side                    string  `json:"side"`
	Price                   float64 `json:"price"`
	Type                    string  `json:"type"`
	Size                    float64 `json:"size"`
	ReduceOnly              bool    `json:"reduceOnly"`
	Ioc                     bool    `json:"ioc"`
	PostOnly                bool    `json:"postOnly"`
	ExternalReferralProgram string  `json:"externalReferralProgram"`
	// ClientID                string  `json:"clientId"`
}

func (f *Client) SetOrder(ticker string, side string, price, size float64, ordertype string, reduceOnly bool) (NewOrderResponse, error) {
	var newOrderResponse NewOrderResponse
	requestBody, err := json.Marshal(NewOrder{
		Market:     ticker,
		Side:       side,
		Price:      price,
		Type:       ordertype,
		Size:       size,
		ReduceOnly: reduceOnly,
		Ioc:        false,
		PostOnly:   reduceOnly})
	if err != nil {
		log.Printf("Error PlaceOrder", err)
		return newOrderResponse, err
	}
	resp, err := f.post("orders", requestBody)
	if err != nil {
		log.Printf("Error PlaceOrder", err)
		return newOrderResponse, err
	}
	err = processResponse(resp, &newOrderResponse)
	return newOrderResponse, err
}
