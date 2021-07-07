package tracking

import "github.com/DawnKosmos/ftxcmd/ftx"

type WithdrawsDeposits struct {
	Deposits  []ftx.Deposit  `json:"deposits,omitempty"`
	Withdraws []ftx.Withdraw `json:"withdraws,omitempty"`
}

func (f *WithdrawsDeposits) AddDeposit(v ftx.Deposit) {
	f.Deposits = append(f.Deposits, v)
}
func (f *WithdrawsDeposits) AddWithdraw(v ftx.Withdraw) {
	f.Withdraws = append(f.Withdraws, v)
}
