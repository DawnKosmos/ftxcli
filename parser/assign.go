package parser

import (
	"errors"
)

/*

ASSIGN ::=   { STRING ~ "="  ~ VARIABLE}


*/

func ParseAssign(name string, tk []Token, ws *WsAccount) (err error) {
	if len(tk) == 0 {
		return errors.New("Nothing can't be assigned to a Variable")
	}
	vll := vl
	if ws != nil {
		vll = ws.vl
	}

	switch tk[0].Type {
	case FUNC:
		r, err := ParseAssignFunc(name, tk[1:])
		if err != nil {
			return err
		}
		vll[name] = Variable{FUNCTION, r}
	case SIDE:
		vll[name] = Variable{EXPRESSION, tk}
	case STOP:
		vll[name] = Variable{EXPRESSION, tk}
	default:
		vll[name] = Variable{CONSTANT, tk}
	}

	if ws != nil {
		ws.vl = vll
	} else {
		vl = vll
	}
	return nil
}

func ParseAssignFunc(name string, tk []Token) (f Func, err error) {
	if len(tk) == 0 { //check Input
		return f, errors.New("Empty Func can't be assigned to a Variable")
	}
	if tk[0].Type != LBRACKET { //check Input
		return f, errors.New("Syntax error, no bracket" + tk[0].Text)
	}

	f.Name = name
	m := make(map[string]int) //A map that track which variable is on which position of the tokenlist

	nl := tk[1:]
	var count int
L: //Read the variable
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

	for i, v := range tk { //Create the unparsed tokenlist with VARIABLE token are used as placeholder
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
