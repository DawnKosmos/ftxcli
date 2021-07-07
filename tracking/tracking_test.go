package tracking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

func TestFileCreation(t *testing.T) {

	if ttt == nil {
		t.Fail()
	}

	data, err := ioutil.ReadFile("main.acc")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	s := strings.Split(string(data), " ")
	c := &http.Client{}

	f := ftx.NewClient(c, s[1], s[2], s[0])

	ttt.Fill(f, time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC), time.Now())

	file, err := os.Create("test.txt")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	bytes, err := json.MarshalIndent(ttt, "", "  ")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

}
