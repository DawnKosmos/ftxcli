package parser

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	c := "y = buy btc-perp 10% -high 3h 1%"
	r, err := Lexer(c)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(r)
	_, err = Parse(r)

	if err != nil {
		t.Fatal(err)
	}

	r, err = Lexer("y")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(r)
	/*
		o, err := ParseOrder("buy", r)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(o)*/
}
