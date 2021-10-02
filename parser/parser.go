package parser

import (
	"errors"

	"github.com/DawnKosmos/ftxcmd/ftx"
	"github.com/gorilla/websocket"
)

/*
TODO:
Go Routines for faster execution
Creating expected result and showing the result of the expression
*/

//Engine
type Engine struct {
	Vl      map[string]Variable
	Account ftx.Client
}

type WsAccount struct {
	ws      *websocket.Conn
	buffer  []byte
	vl      map[string]Variable
	account *ftx.Client
}

func NewWsAccount(ws *websocket.Conn, account *ftx.Client) *WsAccount {
	return &WsAccount{ws: ws, account: account, vl: make(map[string]Variable), buffer: []byte{}}
}

func (w *WsAccount) Write(s string) {
	if s != "" {
		w.AddToBuffer(s)
	}
	w.ws.WriteMessage(1, w.buffer)
	w.buffer = []byte{}
}

func (w *WsAccount) AddToBuffer(s string) {
	w.buffer = append(w.buffer, []byte(s)...)
}

//the vl stands for variable list, it saves the varialbes so that they are parsed into new experssions
var vl map[string]Variable

//Parse the funtion and returning and Evaluater, if a variable gets assign nil will be returned
func Parse(tl []Token, ws *WsAccount) (Evaluater, error) {
	nl := tl

	var vll map[string]Variable
	if ws == nil {
		//init the vl list
		if vl == nil {
			vl = make(map[string]Variable)
			vll = vl
		}
	} else {
		vll = ws.vl
	}

	var err error

	/*
		An expression starts with either a
		- variable thats gets assigned
		- a variable
		- order
	*/
	if tl[0].Type == VARIABLE {
		v, ok := vll[tl[0].Text]
		if !ok {
			if len(tl) == 1 {
				return nil, errors.New("THE VARIABLE IS UNKNOWN " + tl[0].Text)
			}
			if tl[1].Type == ASSIGN {
				err = ParseAssign(tl[0].Text, tl[2:], ws)
				return nil, err
			} else {
				return nil, errors.New("THE VARIABLE IS UNKNOWN " + tl[0].Text)
			}
		}

		if len(tl) > 2 {
			if tl[1].Type == ASSIGN {
				delete(vll, tl[0].Text)
				if ws == nil {
					vl = vll
				} else {
					ws.vl = vll
				}
				err = ParseAssign(tl[0].Text, tl[2:], ws)
				return nil, err
			}
		}
		nl, err = ParseVariable(v, tl[1:])
		if err != nil {
			return nil, err
		}
	}

	switch nl[0].Type {
	case SIDE:
		o, err := ParseOrder(nl[0].Text, nl[1:])
		if err != nil {
			return nil, err
		}
		return o, nil
	case STOP:
		if nl[1].Type == SIDE {
			o, err := ParseStop(nl[1].Text, nl[2:])
			if err != nil {
				return nil, err
			}
			return o, nil
		}
		return nil, errors.New("After A stop a buy/sell has to follow")
	case CANCEL:
		o, err := ParseCancel(nl)
		if err != nil {
			return nil, err
		}
		return o, nil
	case FUNDING:
		o, err := ParseFunding(nl[1:])
		if err != nil {
			return nil, err
		}
		return o, nil
	case FUNDINGRATES:
		o, err := ParseFundingRates(nl[1:])
		if err != nil {
			return nil, err
		}
		return o, nil
	case LOAD:
		o, err := ParseLoad(nl[1:])
		if err != nil {
			return nil, err
		}
		return o, nil
	default:
		return nil, errors.New(nl[0].Type.String() + " Is not a legit command")
	}
}
