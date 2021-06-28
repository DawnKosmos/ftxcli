package parser

import (
	"errors"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

type Engine struct {
	Vl      map[string]Variable
	Account ftx.Client
}

var vl map[string]Variable

type Parser interface {
	Evaluate()
}

func Parse(tl []Token) (Parser, error) {
	nl := tl

	vl = make(map[string]Variable)

	var err error
	if tl[0].Type == VARIABLE {
		v, ok := vl[tl[0].Text]
		if !ok {
			if tl[1].Type == ASSIGN {
				err = ParseAssign(tl[0].Text, tl[2:])
				return nil, err
			} else {
				return nil, errors.New("THE VARIABLE IS UNKNOWN " + tl[0].Text)
			}
		}

		nl, err = ParseVariable(v, tl[1:])
	}

	switch tl[0].Type {
	case SIDE:
		o, err := ParseOrder(tl[0].Text, nl[1:])
		if err != nil {
			return nil, err
		}
		return o, nil
	case STOP:
	default:
		return nil, errors.New(tl[0].Text + " Is not a legit command")
	}

	return nil, nil
}

/*

SIDE ::= {buy | sell}
PRICE ::= {LADERED | FLOAT | PERCENT | DIFF }



func ParseCommand(command string, l []Token) Parser {
	var or Order
	or.Order = command
	if l[0].Type != VARIABLE {
		fmt.Println(l[0], "is not a ticker nor variable")

		v, ok := vl[l[0].Text]
		if ok =

	}

}

VARIABLE
	Constanten
	Functionen



*/
