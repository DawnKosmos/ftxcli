package ftx

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	data, err := ioutil.ReadFile("friend.acc")
	if err != nil {
		fmt.Println("File does not exist", err)
		t.FailNow()
	}

	s := strings.Split(string(data), " ")
	c := &http.Client{}

	f := NewClient(c, s[1], s[2], s[0])

	r, err := f.ChangeLev(10)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(r)
}
