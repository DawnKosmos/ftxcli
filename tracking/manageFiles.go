package tracking

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Tracking struct {
	Username string `json:"username,omitempty"`
	Month    string `json:"month,omitempty"`
	Year     int    `json:"year,omitempty"`
	Date     []Day  `json:"date,omitempty"`
}

func init() {
	t := time.Now()
	y, m, d := t.Date()
	//checkFolder
	if _, err := os.Stat("/log/" + strconv.Itoa(y)); os.IsNotExist(err) {
		err := os.Mkdir("/log/"+strconv.Itoa(y), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	FileName := "/log/" + strconv.Itoa(y) + "/" + m.String() + ".log"
	if _, err := os.Stat(FileName); os.IsNotExist(err) {
		emptyFile, err := os.Create(FileName)
		if err != nil {
			log.Fatal(err)
		}
		var t Tracking

	}

}
