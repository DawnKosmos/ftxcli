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

	acc, _ := reader.ReadString('\n')

	acc = acc[:len(acc)-2]

	data, err := ioutil.ReadFile(acc)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	s := strings.Split(string(data), " ")
	c := &http.Client{}

	f := ftx.NewClient(c, s[1], s[2], s[0])

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		t, err := parser.Lexer(input[:len(input)-1])
		fmt.Println(input)
		fmt.Println(t)
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
			fmt.Println("Variable assigned")
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

func evaluate(s string) {
	//lex :=
}
