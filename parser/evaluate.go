package parser

import "github.com/DawnKosmos/ftxcmd/ftx"

/*
Variable Zuweisung, brauch kein Evaluate
Buy, Sell und Stop brauchen ein Evaluate
sowie die extra optionen, -low -high -l -le


-----
Konzepte wie evaluation funktioniert



Order wird geparsed

kriegt je nach OrderType funktionen geschickt, die ausgeführt werden

Price kriegt input und spielt die funktionen

*/

//Evaluater interface implements the evaluate function, which returns an error, TODO: will also return a context
type Evaluater interface {
	Evaluate(f *ftx.Client) error
}
