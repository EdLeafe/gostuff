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
    ballname := " MegaBall"
    args := os.Args[1:]
    if len(args) > 0 {
        game = strings.ToLower(args[0])
    }
    if strings.HasPrefix(game, "pow") {
        game = "power"
        ballname = "PowerBall"
    }

	res := games.QuickPick(game)
	picks := []int{res.P0, res.P1, res.P2, res.P3, res.P4}
	fmt.Printf("Your numbers are: ")
	for _, num := range picks {
		fmt.Printf("%d ", num)
	}
	fmt.Printf("\n")
	fmt.Printf("       %s: %d\n", ballname, res.Ball)
}
