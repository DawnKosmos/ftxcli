package ftxcmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/DawnKosmos/ftxcmd/ftx"
	"github.com/DawnKosmos/ftxcmd/parser"
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

	fmt.Println(s)
	f := ftx.NewClient(c, s[1], s[2], "")
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
		p, err := parser.Parse(t)
		if p == nil {
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(t[0].Text, "Variable assigned")
			continue
		}

		err = p.Evaluate(f)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("SUCCESS")
	}

}
