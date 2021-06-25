package parser

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	c := "y = buy btc-perp 10% -high 5h -l 5 -10 -30"
	r, err := Lexer(c)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(r)
	_, err = Parse(r)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(vl["y"])
	/*
		o, err := ParseOrder("buy", r)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(o)*/
}
