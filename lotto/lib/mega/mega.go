package mega

import (
	"github.com/EdLeafe/numbers"
)

type MegaResult struct {
	P0   int
	P1   int
	P2   int
	P3   int
	P4   int
	Mega int
}

func QuickPick() MegaResult {
	numbers.SeedRandom()
	res := MegaResult{}
	// Pick 5 numbers from 1 to 75
	picks := numbers.RandSet(1, 75, 5)
	// Add the mega ball
	meganum := numbers.RandRange(1, 25)
	res.P0 = picks[0]
	res.P1 = picks[1]
	res.P2 = picks[2]
	res.P3 = picks[3]
	res.P4 = picks[4]
	res.Mega = meganum
	return res
}
