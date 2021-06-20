package ftxcmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/DawnKosmos/ftxcmd/ftx"
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

	s := strings.Split(string(data), "\n")
	c := &http.Client{}

	ftx.NewClient(c, s[1], s[2], s[0])

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)

	}

}

func evaluate(s string) {
	//lex :=
}
