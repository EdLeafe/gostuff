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

func updateSubject(h mail.Header) []string {
	edexp := regexp.MustCompile(EDPAT)
	sub := h.Get("Subject")
	ret := []string{sub, ""}
	if edexp.MatchString(sub) {
		ret[1] = sub
	}
	return ret
}

func updateFrom(h mail.Header) string {
	sendexp := regexp.MustCompile(SENDPAT)
	from := h.Get("From")
	match := sendexp.FindStringSubmatch(from)
	if len(match) > 0 {
		from = match[len(match)-1]
		//fmap[from] += 1
		return from
	}
	return ""
}

func updateRecip(h mail.Header) []string {
	orig := h.Get("X-Original-To")
	ret := []string{"", ""}
	if len(orig) > 0 {
		//rmap[orig] += 1
		ret[0] = orig
	}
	toAddr := h.Get("To")
	if len(toAddr) > 0 {
		//	rmap[toAddr] += 1
		ret[1] = toAddr
	}
	return ret
}

func processMsg(wg *sync.WaitGroup, txt string, chSubj chan string,
	chEds chan string, chSend chan string, chRecp chan string) {

	defer wg.Done()
	stxt := strings.TrimSpace(txt)
	r := strings.NewReader(stxt)
	msg, err := mail.ReadMessage(r)
	check(err)
	h := msg.Header

	ss := updateSubject(h)
	chSubj <- ss[0]
	chEds <- ss[1]
	f := updateFrom(h)
	chSend <- f
	rr := updateRecip(h)
	chSend <- rr[0]
	chRecp <- rr[1]
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
	chSubj := make(chan string, 10000)
	chEds := make(chan string, 10000)
	chSend := make(chan string, 10000)
	chRecp := make(chan string, 10000)
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
				go processMsg(&wg, txt, chSubj, chEds, chSend, chRecp)
			}
		}
	}
	// We've gotten to the end, so the buffer will contain the last message.
	// Process that one, and we're done.
	wg.Add(1)
	go processMsg(&wg, buf.String(), chSubj, chEds, chSend, chRecp)
	wg.Wait()
    close(chSubj)
    close(chEds)
    close(chSend)
    close(chRecp)
	for ss := range chSubj {
		subjs[ss] += 1
	}
	for ss := range chEds {
		eds[ss] += 1
	}
	for ss := range chSend {
		senders[ss] += 1
	}
	for ss := range chRecp {
		recips[ss] += 1
	}

	sortedSubjects := sortStringCount(subjs, false)
	sortedRecips := sortStringCount(recips, true)
	sortedEds := sortStringCount(eds, false)
	sortedSenders := sortStringCount(senders, false)

	if sortedSubjects == nil || sortedEds == nil || sortedSenders == nil || sortedRecips == nil{
		fmt.Println("Just to keep the compiler from complaining.")
	}
	for _, itm := range sortedSubjectsdasdasdas{
		fmt.Println(itm)
	}
}
