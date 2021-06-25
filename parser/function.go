package parser

import (
	"errors"
	"strconv"
)

type Func struct {
	Name               string
	NumberOfParameters int
	ParameterPosition  []int
	FunctionTokens     []Token
}

func ParseFunc(v interface{}, tl []Token) ([]Token, error) {
	fun, ok := v.(Func)

	if !ok {
		return tl, errors.New("This variable is not a function")
	}
	e, err := fun.Parse(tl)
	if err != nil {
		return tl, err
	}

	return e, nil
}

func (f *Func) Parse(tl []Token) ([]Token, error) {
	if tl[0].Type != LBRACKET {
		return tl, errors.New("Syntax error, no bracket" + tl[0].Text)
	}
	var ss string
	for _, v := range tl[1:] {
		switch v.Type {
		case RBRACKET:
			break
		case VARIABLE:
			ss += v.Text + " "
		default:
			return tl, errors.New("Undentified error")
		}
	}

	s, err := Lexer(ss)

	if err != nil {
		return tl, err
	}

	if len(s) != f.NumberOfParameters {
		return tl, errors.New("To much Parameters, amount has to be" + strconv.Itoa(f.NumberOfParameters))
	}

	for i, t := range s {
		n := f.ParameterPosition[i]
		f.FunctionTokens[n] = t
	}

	return f.FunctionTokens, nil
}
