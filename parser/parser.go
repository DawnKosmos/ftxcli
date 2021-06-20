package parser

import "fmt"

type VariableType int

const (
	FUNCTION VariableType = iota
	EXPRESSION
	NUM
	NAME
	EXCHANGE
)

var vl map[string]Variable

func Init(){
	v1 = make(map[string]Variable)
}

type Variable struct {
	Type    VariableType
	Content interface{}
}

type Parser interface {
	Evaluate()
}

func Parse(l []Token) {
	switch l[0].Type {
	case COMMAND:
	case VARIABLE:
	default:
	}
}

type OrderType int

const (
	SIMPLE OrderType = iota
	CASCADE
)

type Order struct {
	Type   OrderType
	Order  string
	Ticker string
	Size   float64
}

func ParseCommand(command string, l []Token) Parser {
	var or Order
	or.Order = command
	if l[0].Type != VARIABLE {
		fmt.Println(l[0], "is not a ticker nor variable")
		if 
	}

}

/*
VARIABLE
	Constanten
	Functionen




*/
