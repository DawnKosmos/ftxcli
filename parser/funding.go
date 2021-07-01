package parser

/*

funding -position //funding rate der aktuellen positionen
funding -highest 20 //funding der highest 20 coins
funding 10
funding 6h
funding -sum 1d
funding btc eth 10
*/

type FundingType int

const (
	PAYMENTS FundingType = iota
	POSITIONS
	GENERAL
)

type Funding struct {
	ft        FundingType
	Ticker    string
	Time      int64
	Summarize bool
}

func ParseFunding(tl []Token) (*Funding, error) {

}
