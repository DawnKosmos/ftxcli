package tracking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

var ttt *Tracker

func init() {
	t := time.Now()
	y, m, _ := t.Date()
	//checkFolder
	var err error

	if _, err = os.Stat("log/"); os.IsNotExist(err) {
		err := os.Mkdir("log/", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err = os.Stat("log/" + strconv.Itoa(y)); os.IsNotExist(err) {

		err := os.Mkdir("log/"+strconv.Itoa(y), 0755)

		if err != nil {
			log.Fatal(err)
		}
	}

	FileName := "log/" + strconv.Itoa(y) + "/" + m.String() + ".log"
	if _, err := os.Stat(FileName); os.IsNotExist(err) {
		emptyFile, err := os.Create(FileName)
		defer emptyFile.Close()
		if err != nil {
			log.Fatal(err)
		}
		var t Tracker
		err = FillFile(emptyFile, m, &t)
		if err != nil {
			fmt.Println(err)
		}
	}
	//Load month

	ttt = &Tracker{}

	data, err := ioutil.ReadFile(FileName)
	if err != nil {
		fmt.Println("File to track account haven't been found")
		return
	}

	err = json.Unmarshal(data, ttt)
	if err != nil {
		fmt.Println("File to track account haven't been found")
		return
	}

}

func FillFile(f *os.File, m time.Month, t *Tracker) error {
	t.Month = m.String()
	for i := 1; i < 32; i++ {
		t.Date = append(t.Date, Day{Day: i})
	}

	res, err := json.MarshalIndent(*t, "", "\t")
	if err != nil {
		return err
	}

	_, err = f.WriteString(string(res))
	if err != nil {
		return err
	}

	return nil
}
