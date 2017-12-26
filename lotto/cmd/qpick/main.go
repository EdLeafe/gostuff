package main

import (
	"fmt"
    "github.com/EdLeafe/numbers"
)

func main() {
    numbers.SeedRandom()
    // Pick 5 numbers from 1 to 75
    picks := numbers.RandSet(1, 75, 5)

    // Add the mega ball
    mega := numbers.RandRange(1, 25)

    fmt.Printf("Your numbers are: ")
    for _, num := range picks {
        fmt.Printf("%d ", num)
    }
    fmt.Printf("\n")
    fmt.Printf("            Mega: %d\n", mega)
}
