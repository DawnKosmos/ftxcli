package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/DawnKosmos/ftxcmd/ftx"
)

type Load struct {
	FileName []string
}

func ParseLoad(tl []Token) (*Load, error) {
	if len(tl) == 0 {
		return nil, errors.New("A file has to be provided to load")
	}
	l := &Load{FileName: []string{}}

	for _, v := range tl {
		switch v.Type {
		case VARIABLE:
			l.FileName = append(l.FileName, v.Text)
		case FLAG:
			continue
		default:
			fmt.Println(v.Type.String(), v.Text, "is not a good filename")
			continue
		}
	}

	return l, nil
}

func (l *Load) Evaluate(f *ftx.Client) error {
	if len(l.FileName) == 0 {
		return errors.New("no valid filenames")
	}
	var data []byte
	var err error

	for i, v := range l.FileName {
		data, err = ioutil.ReadFile(v)
		if err != nil {
			fmt.Println(v, err)
			continue
		}

		s := strings.Split(string(data), "\n")
		for j, vv := range s {
			if vv[len(vv)-1] == byte(13) {
				vv = vv[:len(vv)-1]
			}

			ss, err := Lexer(vv)
			if err != nil {
				fmt.Println("file", i, " line", j, err)
			}
			if ss[0].Type != VARIABLE {
				continue
			}

			parsed, err := Parse(ss)
			if parsed != nil {
				fmt.Println("file", i, " line", j, " not assignmed")
			}
			if err != nil {
				fmt.Println("file", i, " line", j, err)
			}
			fmt.Println(ss[0].Text, " Variable assigned")
		}
	}
	return nil
}
