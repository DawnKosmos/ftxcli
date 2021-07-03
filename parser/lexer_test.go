package parser

import (
	"fmt"
	"testing"
)

/*
func TestAmount(t *testing.T) {

	data, err := ioutil.ReadFile("main.acc")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	s := strings.Split(string(data), " ")
	c := &http.Client{}

	f := ftx.NewClient(c, s[1], s[2], s[0])
	a := Amount{
		Type: ACCOUNTSIZE,
		Val:  200,
	}

	v, err := a.Evaluate(f, "btc-perp")

	fmt.Println(v)
}*/

func TestLadder(t *testing.T) {
	/*
		data, err := ioutil.ReadFile("main.acc")
		if err != nil {
			fmt.Println("File reading error", err)
			return
		}

		s := strings.Split(string(data), " ")
		c := &http.Client{}

		f := ftx.NewClient(c, s[1], s[2], s[0])
	*/

	s := "x(u100,-100,-200)"

	o, err := Lexer(s)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(o)
}
