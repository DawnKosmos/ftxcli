package parser

import "errors"

type VariableType int

const (
	FUNCTION VariableType = iota
	EXPRESSION
	CONSTANT
)

type Variable struct {
	Type    VariableType
	Content interface{}
}

func ParseVariable(v Variable, tl []Token) ([]Token, error) {
	switch v.Type {
	case FUNCTION:
		l, err := ParseFunc(v.Content, tl[1:])
		if err != nil {
			return tl, err
		}
		return l, nil

	case EXPRESSION:
		o, ok := v.Content.([]Token)
		if !ok {
			return tl, errors.New("Some error occured on my side, sry i guess")
		}
		return o, nil
	case CONSTANT:
		o, ok := v.Content.([]Token)
		if !ok {
			return tl, errors.New("Some error occured on my side, sry i guess")
		}
		return o, nil
	}

	return nil, errors.New("Something went absolute wrong with me, sry lol")
}
