package parser

import (
	"errors"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

/*
TODO:
Cancle, Funding, Stop
Go Routines for faster execution
Creating expected result and showing the result of the expression
*/

//Engine
type Engine struct {
	Vl      map[string]Variable
	Account ftx.Client
}

//the vl stands for variable list, it saves the varialbes so that they are parsed into new experssions
var vl map[string]Variable

//Parse the funtion and returning and Evaluater, if a variable gets assign nil will be returned
func Parse(tl []Token) (Evaluater, error) {
	nl := tl

	//init the vl list
	if vl == nil {
		vl = make(map[string]Variable)
	}
	var err error

	/*
		An expression starts with either a
		- variable thats gets assigned
		- a variable
		- order
	*/
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

		if len(tl) > 2 {

			if tl[1].Type == ASSIGN {
				delete(vl, tl[0].Text)
				err = ParseAssign(tl[0].Text, tl[2:])
				return nil, err
			}
		}
		nl, err = ParseVariable(v, tl[1:])
	}

	switch nl[0].Type {
	case SIDE:
		o, err := ParseOrder(nl[0].Text, nl[1:])
		if err != nil {
			return nil, err
		}
		return o, nil
	case STOP:
	case CANCEL:
		o, err := ParseCancel(nl)
		if err != nil {
			return nil, err
		}
		return o, nil
	case FUNDING:
	default:
		return nil, errors.New(nl[0].Text + " Is not a legit command")
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
