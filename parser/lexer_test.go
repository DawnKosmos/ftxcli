package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

func TestAmount(t *testing.T) {

	data, err := ioutil.ReadFile("main.acc")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	s := strings.Split(string(data), " ")
	c := &http.Client{}

	f := ftx.NewClient(c, s[1], s[2], s[0])

	b, err := f.GetFills()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	for _, v := range b {
		fmt.Println(v)

	}
}

/*
func TestLadder(t *testing.T) {
		data, err := ioutil.ReadFile("main.acc")
		if err != nil {
			fmt.Println("File reading error", err)
			return
		}

		s := strings.Split(string(data), " ")
		c := &http.Client{}

		f := ftx.NewClient(c, s[1], s[2], s[0])

		s := "load variables.acc"
		l, err := Lexer(s)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}

		p, err := Parse(l)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}

		err = p.Evaluate(nil)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}

		/*for k, v := range vl {
			fmt.Println(k, v)

		}

*/
