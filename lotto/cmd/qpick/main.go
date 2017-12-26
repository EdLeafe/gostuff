package main

import (
	"fmt"
	"github.com/EdLeafe/lotto/lib/mega"
)

func main() {
	res := mega.QuickPick()
	picks := []int{res.P0, res.P1, res.P2, res.P3, res.P4}
	fmt.Printf("Your numbers are: ")
	for _, num := range picks {
		fmt.Printf("%d ", num)
	}
	fmt.Printf("\n")
	fmt.Printf("            Mega: %d\n", res.Mega)
}
