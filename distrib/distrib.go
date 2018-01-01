package main

import (
	"fmt"
    "os"
	"github.com/EdLeafe/numbers"
    "strconv"
    "text/tabwriter"
)

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
        panic("You need to supply low, high, and total")
    }
    slow := args[0]
    shigh := args[1]
    stotal := args[2]
    low, _ := strconv.Atoi(slow)
    high, _ := strconv.Atoi(shigh)
    total, _ := strconv.Atoi(stotal)
    numbers.SeedRandom()
    diff := high - low + 1
    count := make([]int, diff)

    rr := 0
    for i:=0; i < total; i++ {
        rr = numbers.RandRange(low, high)
        count[rr - low] += 1
    }

    out := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0) //tabwriter.AlignRight) // | tabwriter.Debug)
    txt := "Number\tOccurs\t"
    fmt.Fprintln(out, txt)

    for i:=low; i<= high; i++ {
        txt = fmt.Sprintf("%v\t%v\t", i, count[i - low])
        fmt.Fprintln(out, txt)
    }
    out.Flush()
}

