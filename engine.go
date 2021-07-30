package ftxcmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/DawnKosmos/ftxcmd/ftx"
	"github.com/DawnKosmos/ftxcmd/parser"
	"github.com/gorilla/websocket"
)

func Start(e ...interface{}) {
	fmt.Print("Welcome, write name of acc file to import the keys\n>")
	reader := bufio.NewReader(os.Stdin)
ReadAcc:
	acc, _ := reader.ReadString('\n')
	acc = acc[:len(acc)-2]
	data, err := ioutil.ReadFile(acc)
	if err != nil {
		fmt.Println("File does not exist", err)
		goto ReadAcc
	}

	s := strings.Split(string(data), " ")
	c := &http.Client{}

	if s[0] == "main" {
		s[0] = ""
	}

	f := ftx.NewClient(c, s[1], s[2], s[0])
	a, err := f.GetAccount()
	if err != nil {
		fmt.Println("Account Verification failed", err)
		return
	}

	fmt.Println(a.Username, " Account Value:", a.TotalAccountValue, " Amount of Positions:", len(a.Positions))

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		t, err := parser.Lexer(input[:len(input)-1])
		if err != nil {
			fmt.Println(err)
			continue
		}
		p, err := parser.Parse(t, nil)
		if p == nil {
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(t[0].Text, "Variable assigned")
			continue
		}

		err = p.Evaluate(f, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("SUCCESS")
	}

}

func WsStart(wsc *websocket.Conn, name, public, secret string) {

	f := ftx.NewClient(&http.Client{}, public, secret, name)
	a, err := f.GetAccount()
	if err != nil {
		wsc.WriteMessage(1, []byte("Account Verification failed"))
		return
	}

	ws := parser.NewWsAccount(wsc, f)

	b := strconv.Itoa(int(a.TotalAccountValue))
	wsc.WriteMessage(1, []byte(a.Username+" Account Value: "+b+" Amount of Positions: "+strconv.Itoa(len(a.Positions))))

	for {
		_, msg, err := wsc.ReadMessage()
		if err != nil {
			fmt.Println(err)
		}

		t, err := parser.Lexer(string(msg))
		if err != nil {
			fmt.Println(err)
			continue
		}
		p, err := parser.Parse(t, ws)
		if p == nil {
			if err != nil {
				fmt.Println(err)
				continue
			}
			ws.Write("Variable assigned: " + t[0].Text)
			continue
		}

		err = p.Evaluate(f, ws)
		if err != nil {
			ws.Write(err.Error())
			continue
		}
		fmt.Println("SUCCESS")
	}
}
