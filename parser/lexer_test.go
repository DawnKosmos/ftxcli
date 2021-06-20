package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	c := "buy 10% -l 38030 38070 5"

	r, _ := Lexer(strings.Split(c, " "))
	for _, v := range r {
		fmt.Println(v)
	}
}
