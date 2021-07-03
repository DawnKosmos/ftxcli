package parser

import (
	"errors"
	"fmt"
)

/*

ASSIGN ::=   { STRING ~ "="  ~ VARIABLE}


*/

func ParseAssign(name string, tk []Token) (err error) {
	if len(tk) == 0 {
		return errors.New("Nothing can't be assigned to a Variable")
	}

	switch tk[0].Type {
	case FUNC:
		r, err := ParseAssignFunc(name, tk[1:])
		if err != nil {
			return err
		}
		vl[name] = Variable{FUNCTION, r}
	case SIDE:
		vl[name] = Variable{EXPRESSION, tk}
	case STOP:
		vl[name] = Variable{EXPRESSION, tk}
	default:
		fmt.Println(tk[0].Type)
		vl[name] = Variable{CONSTANT, tk}
	}

	return nil
}

func ParseAssignFunc(name string, tk []Token) (f Func, err error) {
	if tk[0].Type != LBRACKET {
		return f, errors.New("Syntax error, no bracket" + tk[0].Text)
	}
	f.Name = name
	m := make(map[string]int)

	nl := tk[1:]
	var count int
L:
	for _, v := range nl {
		switch v.Type {
		case RBRACKET:
			break L
		case VARIABLE:
			m[v.Text] = count
		default:
			return f, errors.New(v.Text + " is not a valid parameter name")
		}

		count++
	}

	f.ParameterPosition = make([]int, count)

	tk = tk[count+2:]

	f.NumberOfParameters = count

	for i, v := range tk {
		if v.Type == VARIABLE {
			n, ok := m[v.Text]
			if ok {
				f.ParameterPosition[n] = i
			}
		}
		f.FunctionTokens = append(f.FunctionTokens, v)
	}

	return f, nil
}
