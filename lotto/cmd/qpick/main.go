package main

import (
	"fmt"
	"github.com/EdLeafe/lotto/lib/games"
    "os"
    "strings"
)

func main() {
    // Default to mega millions
    game := "mega"
    range1 := 75
    range2 := 25
    ballname := " Mega"
    args := os.Args[1:]
    if len(args) > 0 {
        game = strings.ToLower(args[0])
    }
    if strings.HasPrefix(game, "pow") {
        range1 = 69
        range2 = 26
        ballname = "Power"
    }

	res := games.QuickPick(range1, range2)
	picks := []int{res.P0, res.P1, res.P2, res.P3, res.P4}
	fmt.Printf("Your numbers are: ")
	for _, num := range picks {
		fmt.Printf("%d ", num)
	}
	fmt.Printf("\n")
	fmt.Printf("           %s: %d\n", ballname, res.Ball)
}
