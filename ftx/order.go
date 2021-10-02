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

type NewTriggerOrder struct {
	Market       string  `json:"market"`
	Side         string  `json:"side"`
	Size         float64 `json:"size"`
	Type         string  `json:"type"`
	ReduceOnly   bool    `json:"reduceOnly"`
	TriggerPrice float64 `json:"triggerPrice"`
}

func (f *Client) SetOrder(ticker string, side string, price, size float64, orderType string, reduceOnly bool) (NewOrderResponse, error) {
	var newOrderResponse NewOrderResponse
	requestBody, err := json.Marshal(NewOrder{
		Market:     ticker,
		Side:       side,
		Price:      price,
		Type:       orderType,
		Size:       size,
		ReduceOnly: reduceOnly,
		Ioc:        false,
		PostOnly:   reduceOnly})
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

type NewTriggerOrderResponse struct {
	Success bool         `json:"success"`
	Result  TriggerOrder `json:"result"`
}

type TriggerOrder struct {
	CreatedAt    time.Time `json:"createdAt"`
	Error        string    `json:"error"`
	Future       string    `json:"future"`
	ID           int64     `json:"id"`
	Market       string    `json:"market"`
	OrderID      int64     `json:"orderId"`
	OrderPrice   float64   `json:"orderPrice"`
	ReduceOnly   bool      `json:"reduceOnly"`
	Side         string    `json:"side"`
	Size         float64   `json:"size"`
	TrailStart   float64   `json:"trailStart"`
	TrailValue   float64   `json:"trailValue"`
	TriggerPrice float64   `json:"triggerPrice"`
	TriggeredAt  string    `json:"triggeredAt"`
	Type         string    `json:"type"`
	OrderType    string    `json:"orderType"`
	FilledSize   float64   `json:"filledSize"`
	AvgFillPrice float64   `json:"avgFillPrice"`
	OrderStatus  string    `json:"orderStatus"`
}

func (f *Client) SetTriggerOrder(ticker string, side string, price, size float64, orderType string, reduceOnly bool) (NewTriggerOrderResponse, error) {
	var newOrderResponse NewTriggerOrderResponse
	requestBody, err := json.Marshal(NewTriggerOrder{
		Market:       ticker,
		Side:         side,
		TriggerPrice: price,
		Size:         size,
		Type:         orderType,
		ReduceOnly:   reduceOnly})
	if err != nil {
		log.Printf("Error PlaceOrder %v", err)
		return newOrderResponse, err
	}
	resp, err := f.post("conditional_orders", requestBody)

	if err != nil {
		log.Printf("Error PlaceTriggerOrder %v", err)
		return newOrderResponse, err
	}
	err = processResponse(resp, &newOrderResponse)
	return newOrderResponse, err
}

type CancelOrders struct {
	Market          string `json:"market,omitempty"`
	ConditionalOnly bool   `json:"conditionalOrdersOnly,omitempty"`
	LimitOnly       bool   `json:"limitOrdersOnly,omitempty"`
}

func (f *Client) DeleteOrders(market string, conditionalOnly, limitOnly bool) (Response, error) {
	var resp Response
	requestBody, err := json.Marshal(CancelOrders{
		Market:          market,
		ConditionalOnly: conditionalOnly,
		LimitOnly:       limitOnly,
	})
	if err != nil {
		log.Printf("Error Delete Orders  %v", err)
		return resp, err
	}
	respRequest, err := f.delete("orders", requestBody)

	if err != nil {
		log.Printf("Error Delete Orders %v", err)
		return resp, err
	}
	err = processResponse(respRequest, &resp)
	return resp, err
}

func (f *Client) CancelOrders(market string, Side string, trigger bool) error {
	var resp Response
	requestBody, err := json.Marshal(Corder{
		Market:               market,
		Side:                 Side,
		ConditionlOrdersOnly: trigger,
	})
	if err != nil {
		log.Printf("Error Delete Orders  %v", err)
		return err
	}
	respRequest, err := f.delete("orders", requestBody)

	if err != nil {
		log.Printf("Error Delete Orders %v", err)
		return err
	}
	err = processResponse(respRequest, &resp)
	return err
}

type Corder struct {
	Market               string `json:"market,omitempty"`
	Side                 string `json:"side,omitempty"`
	ConditionlOrdersOnly bool   `json:"conditionlOrdersOnly,omitempty"`
}
