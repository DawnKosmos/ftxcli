package parser

import "github.com/DawnKosmos/ftxcmd/ftx"

//Evaluater interface implements the evaluate function, which returns an error, TODO: will also return a context
type Evaluater interface {
	Evaluate(f *ftx.Client) error
}
