package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/DawnKosmos/ftxcmd/ftx"
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

	data, err := ioutil.ReadFile("main.acc")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	s := strings.Split(string(data), " ")
	c := &http.Client{}

	f := ftx.NewClient(c, s[1], s[2], s[0])
	p := &Price{Type: PERCENTPRICE,
		PC:         "low",
		Duration:   3600,
		IsLaddered: [2]bool{true, true},
		Values:     [3]float64{4, 5, 10},
	}

	err = p.Evaluate(f, "buy", "xrp-perp", 100)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

/*
5 30000 31000 32000 33000 34000

*/
