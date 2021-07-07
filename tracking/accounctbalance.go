package tracking

import (
	"time"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

type AccountBalance struct {
	Data []AccountSnapshot `json:"data,omitempty"`
}

func (f *AccountBalance) Add(v []ftx.Coin) {
	var sum float64
	for _, vv := range v {
		sum += vv.UsdValue
	}
	f.Data = append(f.Data, AccountSnapshot{Total: sum, Coins: v, Time: time.Now()})
}

type AccountSnapshot struct {
	Total float64    `json:"total,omitempty"`
	Coins []ftx.Coin `json:"coin,omitempty"`
	Time  time.Time  `json:"time,omitempty"`
}
