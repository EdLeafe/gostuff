package main

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
//	"net/smtp"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const FNAME = "spammail"
const MSGSPLIT = "\\n\\nFrom "
const EDPAT = "\\bED\\b"
const SENDPAT = "([^<]+)<\\S+>"

type StringCount struct {
	Text  string
	Count int
}
type byStringCount []StringCount

func (sc byStringCount) Len() int {
	return len(sc)
}
func (sc byStringCount) Swap(i, j int) {
	sc[i], sc[j] = sc[j], sc[i]
}
func (sc byStringCount) Less(i, j int) bool {
	// We always need this sorted in reverse order on Count, and increasing on
	// Text.
	if sc[i].Count > sc[j].Count {
		return true
	}
	if sc[i].Count < sc[j].Count {
		return false
	}
	return sc[i].Text < sc[j].Text
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func updateSubject(h mail.Header, smap, emap map[string]int) {
	edexp := regexp.MustCompile(EDPAT)
	sub := h.Get("Subject")
	smap[sub] += 1
	if edexp.MatchString(sub) {
		emap[sub] += 1
	}
}

func updateFrom(h mail.Header, fmap map[string]int) {
	sendexp := regexp.MustCompile(SENDPAT)
	from := h.Get("From")
	match := sendexp.FindStringSubmatch(from)
	if len(match) > 0 {
		from = match[len(match)-1]
		fmap[from] += 1
	}
}

func updateRecip(h mail.Header, rmap map[string]int) {
	orig := h.Get("X-Original-To")
	if len(orig) > 0 {
		rmap[orig] += 1
		return
	}
	toAddr := h.Get("To")
	if len(toAddr) > 0 {
		rmap[toAddr] += 1
	}
}

func processMsg(wg *sync.WaitGroup, txt string, subjs, eds, senders,
        recips map[string]int) {
    defer wg.Done()
	stxt := strings.TrimSpace(txt)
	r := strings.NewReader(stxt)
	msg, err := mail.ReadMessage(r)
	check(err)
	h := msg.Header

	updateSubject(h, subjs, eds)
	updateFrom(h, senders)
	updateRecip(h, recips)
}

func sortStringCount(txts map[string]int, ones bool) []string {
	bsc := byStringCount{}
	for k, v := range txts {
		sc := StringCount{k, v}
		bsc = append(bsc, sc)
	}
	sort.Sort(bsc)

	numtxt := []string{}
	for _, scOrd := range bsc {
		txt := scOrd.Text
		if ones || scOrd.Count > 1 {
			txt = fmt.Sprintf("[%d] %s", scOrd.Count, scOrd.Text)
		}
		numtxt = append(numtxt, txt)
	}
	return numtxt
}

func main() {
	// Create the maps for the various reports
	subjs := map[string]int{}
	eds := map[string]int{}
	senders := map[string]int{}
	recips := map[string]int{}

	var wg sync.WaitGroup
	var buf bytes.Buffer
	splitter := regexp.MustCompile(MSGSPLIT)
	f, err := os.Open(FNAME)
	check(err)
	defer f.Close()
	b := make([]byte, 1024)
	for {
		// Read 1K chunk
		if _, err = f.Read(b); err == io.EOF {
			break
		}
		// Write it to the buffer
		buf.WriteString(string(b))
		// See if the buffer contains the separator
		found := splitter.FindIndex(buf.Bytes()[1:])
		if found != nil {
			start := found[0]
			if start != 0 {
				alltext := buf.String()
				buf.Reset()
				// Account for the offset in the FindIndex call
				pos := start + 1
				txt := alltext[:pos]
				buf.WriteString(alltext[pos:])
                wg.Add(1)
				go processMsg(&wg, txt, subjs, eds, senders, recips)
			}
		}
	}
	// We've gotten to the end, so the buffer will contain the last message.
	// Process that one, and we're done.
	go processMsg(&wg, buf.String(), subjs, eds, senders, recips)
    wg.Wait()

	sortedSubjects := sortStringCount(subjs, false)
	sortedRecips := sortStringCount(recips, true)
	sortedEds := sortStringCount(eds, false)

	if sortedSubjects == nil || sortedRecips == nil {
		fmt.Println("Just to keep the compiler from complaining.")
	}
	for _, itm := range sortedEds {
		fmt.Println(itm)
	}
}
