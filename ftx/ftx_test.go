package ftx

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNewOrder(t *testing.T) {
	c := NewClient(&http.Client{}, "EU5AinoxN4LnkYlmP1fJb9E2xwTj4pWawKyIuejS", "F7QkEVvg6mrg7RIw9K-GaUXR10Yxi90aZjH27yBJ", "100to1000")
	/*	r, err := c.SetOrder("xrp-perp", "buy", 10, 10, "market", false)

		if err != nil {
			t.Errorf("Error %v", err)
		}
		fmt.Println(r.Result)*/

	b, err := c.SetTriggerOrder("xrp-perp", "sell", 0.50, 10, "stop", false)

	if err != nil {
		t.Errorf("Error %v", err)
	}
	fmt.Println(b.Result)
}

/*
100to1000
EU5AinoxN4LnkYlmP1fJb9E2xwTj4pWawKyIuejS
F7QkEVvg6mrg7RIw9K-GaUXR10Yxi90aZjH27yBJ
*/
