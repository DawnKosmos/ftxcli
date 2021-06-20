package parser

import (
	"fmt"
	"strconv"
)

/*
Syntax: buy, sell, stop, %, positions, open orders, ladder -e 2%
1]

buy btc-perp [0.1, u100, 10%] [market, 0.1%, 38112] -s [1%, -100, 36800]

[-l, -le] [1%, -100 -300, 37100 37400]

x = func(p) buy

x(btc-perp)

*/

type TokenType int

const (
	EOL TokenType = iota
	COMMAND
	FLAG
	NUMBER
	PERCENT
	VARIABLE
	ASSIGN
	FUNC
	LBRACKET
	RBRACKET
)

type Token struct {
	Type TokenType
	Text string
}

func Lexer(input []string) (t []Token, err error) {

	for _, s := range input {
		switch s {
		case "buy", "Buy":
			t = append(t, Token{COMMAND, "buy"})
		case "sell", "Sell":
			t = append(t, Token{COMMAND, "sell"})
		case "=":
			t = append(t, Token{ASSIGN, "="})
		case "func":
			t = append(t, Token{FUNC, "func"})
		default:
			if s[0] == '-' {
				t = append(t, Token{FLAG, s})
				continue
			}
			if s[len(s)-1] == '%' {
				_, err := strconv.Atoi(s[:len(s)-1])
				if err != nil {
					return t, err
				}
				t = append(t, Token{PERCENT, s[:len(s)-1]})
				continue
			}

			_, err := strconv.ParseFloat(s, 64)
			if err == nil {
				t = append(t, Token{NUMBER, s})
				continue
			}

			t = append(t, Token{VARIABLE, s})
		}
	}
	fmt.Println("kek")
	return t, nil
}

func EvaluateVariable(s []byte) []Token {
	var temp []byte
	var tk []Token

	for _, v := range s {
		switch v {
		case '(':
			tk = append(tk, Token{VARIABLE, string(temp)}, Token{LBRACKET, ""})
			temp = []byte("")
		case ')':
			tk = append(tk, Token{VARIABLE, string(temp)}, Token{RBRACKET, ""})
			temp = []byte("")
		case ',':
			tk = append(tk, Token{VARIABLE, string(temp)})
			temp = []byte("")
		default:
			temp = append(temp, v)
		}
	}

	return tk
}
