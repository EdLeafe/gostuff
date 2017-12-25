package main

import (
	"fmt"
    "github.com/EdLeafe/numbers"
    "github.com/EdLeafe/arrayFunc"
	"os"
	"strconv"
    "text/tabwriter"
    "time"
)

type matchResult struct {
    match0 int
    match1 int
    match2 int
    match3 int
    match4 int
    match5 int
    match0M int
    match1M int
    match2M int
    match3M int
    match4M int
    match5M int
}

type megaResult struct {
    one int
    two int
    three int
    four int
    five int
    mega int
}

type matchCompare struct {
    number int
    mega bool
}

func MegaMillionsPick() (megaResult) {
        res := megaResult{}
		picks := numbers.RandSet(1, 70, 5)
		res.one = picks[0]
        res.two = picks[1]
        res.three = picks[2]
        res.four = picks[3]
        res.five = picks[4]
		res.mega = numbers.RandRange(1, 25)
        return res
}

func matchDrawing(ticket megaResult, drawing megaResult) matchCompare {
    megaMatch := ticket.mega == drawing.mega
    ticketArray := []int{ticket.one, ticket.two, ticket.three, ticket.four,
            ticket.five}
    drawingArray := []int{drawing.one, drawing.two, drawing.three, drawing.four,
            drawing.five}
    matches := 0
    for val, _ := range ticketArray {
        if arrayFunc.IntIn(val, drawingArray) {
            matches += 1
        }
    }
    return matchCompare{matches, megaMatch}
}

func addToResults(meta matchCompare, results *matchResult) {
    switch meta {
    case matchCompare{0, false}:
        results.match0 += 1
    case matchCompare{1, false}:
        results.match1 += 1
    case matchCompare{2, false}:
        results.match2 += 1
    case matchCompare{3, false}:
        results.match3 += 1
    case matchCompare{4, false}:
        results.match4 += 1
    case matchCompare{5, false}:
        results.match5 += 1
    case matchCompare{0, true}:
        results.match0M += 1
    case matchCompare{1, true}:
        results.match1M += 1
    case matchCompare{2, true}:
        results.match2M += 1
    case matchCompare{3, true}:
        results.match3M += 1
    case matchCompare{4, true}:
        results.match4M += 1
    case matchCompare{5, true}:
        results.match5M += 1
    }
}

func runDrawings(count int, ticket megaResult, c chan matchCompare) {
    draw := megaResult{}
    for i:=0; i<count; i++ {
        draw = MegaMillionsPick()
        c <- matchDrawing(ticket, draw)
    }
    close(c)
}

func col(label string, num int, pct float64) string {
    return fmt.Sprintf("%v\t%v\t%6.4f\t", label, num, pct)
}

func output(total int, results matchResult, elapsed time.Duration) {
    out := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.AlignRight) // | tabwriter.Debug)
    txt := fmt.Sprintf("The simulation was run for %v drawings", total)
    fmt.Fprintln(out, txt)
    txt = fmt.Sprintf("It took %v to run", elapsed)
    fmt.Fprintln(out, txt)
    fmt.Println("")
    fmt.Fprintln(out, "Count\tResult\t% of draws\t")
    fmt.Fprintln(out, "-----\t------\t----------\t")

    num := results.match0
    pct := (float64(results.match0) / float64(total)) * 100
    fmt.Fprintln(out, col("0", num, pct))
    num = results.match0M
    pct = (float64(results.match0M) / float64(total)) * 100
    fmt.Fprintln(out, col("0 + mega", num, pct))

    num = results.match1
    pct = (float64(results.match1) / float64(total)) * 100
    fmt.Fprintln(out, col("1", num, pct))
    num = results.match1M
    pct = (float64(results.match1M) / float64(total)) * 100
    txt = fmt.Sprintf("1 + mega\t%v\t%6.4f\t", num, pct)
    fmt.Fprintln(out, txt)

    num = results.match2
    pct = (float64(results.match2) / float64(total)) * 100
    fmt.Fprintln(out, col("2", num, pct))
    num = results.match2M
    pct = (float64(results.match2M) / float64(total)) * 100
    fmt.Fprintln(out, col("2 + mega", num, pct))

    num = results.match3
    pct = (float64(results.match3) / float64(total)) * 100
    fmt.Fprintln(out, col("3", num, pct))
    num = results.match3M
    pct = (float64(results.match3M) / float64(total)) * 100
    fmt.Fprintln(out, col("3 + mega", num, pct))

    num = results.match4
    pct = (float64(results.match4) / float64(total)) * 100
    fmt.Fprintln(out, col("4", num, pct))
    num = results.match4M
    pct = (float64(results.match4M) / float64(total)) * 100
    fmt.Fprintln(out, col("4 + mega", num, pct))

    num = results.match5
    pct = (float64(results.match5) / float64(total)) * 100
    fmt.Fprintln(out, col("5", num, pct))
    num = results.match5M
    pct = (float64(results.match5M) / float64(total)) * 100
    fmt.Fprintln(out, col("5 + mega", num, pct))

    out.Flush()
}

func main() {
	// Set the default
    drawings := 1000
	// Check for command-line override
	args := os.Args[1:]
	if len(args) > 0 {
		sdrawings := args[0]
		drawings, _ = strconv.Atoi(sdrawings)
	}
    // Create the struct for the results
    myResults := matchResult{}
    c := make(chan matchCompare, 10)

    // Seed the random engine
    numbers.SeedRandom()

    // OK, let's go!
    start := time.Now()

    myTicket := MegaMillionsPick()
    go runDrawings(drawings, myTicket, c)

    for res := range c {
        addToResults(res, &myResults)
    }
    elapsed := time.Since(start)

    output(drawings, myResults, elapsed)
}
