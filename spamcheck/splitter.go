package main

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

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

func updateRecip(h mail.Header) string {
	orig := h.Get("X-Original-To")
    ret := ""
	if len(orig) > 0 {
		//rmap[orig] += 1
		ret = orig
	}
	return ret
}

func processMsg(wg *sync.WaitGroup, txt string, chans ParseChannels) {
	defer wg.Done()
	stxt := strings.TrimSpace(txt)
	r := strings.NewReader(stxt)
	msg, err := mail.ReadMessage(r)
	check(err)
	h := msg.Header

	ss := updateSubject(h)
	chans.Subjs <- ss[0]
    if ss[1] != "" {
        chans.Eds <- ss[1]
    }
	f := updateFrom(h)
	chans.Senders <- f
	rr := updateRecip(h)
    chans.Recips <- rr
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

type ParseChannels struct {
    Subjs chan string
    Eds chan string
    Senders chan string
    Recips chan string
}

func (ch ParseChannels) Close() {
    close(ch.Subjs)
    close(ch.Eds)
    close(ch.Senders)
    close(ch.Recips)
}

type SpamResults struct {
	SortedSubjects []string
	SortedEds []string
	SortedSenders []string
	SortedRecips []string
    Count int
}

func Analyze(pth string) SpamResults {
    chans := ParseChannels{}
	var wg sync.WaitGroup
	var buf bytes.Buffer
	chans.Subjs = make(chan string, 10000)
	chans.Eds = make(chan string, 10000)
	chans.Senders = make(chan string, 10000)
	chans.Recips = make(chan string, 10000)
	splitter := regexp.MustCompile(MSGSPLIT)
	f, err := os.Open(pth)
	check(err)
	defer f.Close()
	b := make([]byte, 1024)
    msgCount := 0
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
                msgCount += 1
				go processMsg(&wg, txt, chans)
			}
		}
	}
	// We've gotten to the end, so the buffer will contain the last message.
	// Process that one, and we're done.
	wg.Add(1)
    msgCount += 1
	go processMsg(&wg, buf.String(), chans)
	wg.Wait()
    chans.Close()
	// Create the maps for the various reports
    subjs := make(map[string]int)
    eds := make(map[string]int)
    senders := make(map[string]int)
    recips := make(map[string]int)
	for ss := range chans.Subjs {
		subjs[ss] += 1
	}
	for ss := range chans.Eds {
		eds[ss] += 1
	}
	for ss := range chans.Senders {
		senders[ss] += 1
	}
	for ss := range chans.Recips {
		recips[ss] += 1
	}

	sortedSubjs := sortStringCount(subjs, false)
	sortedEds := sortStringCount(eds, false)
	sortedSenders := sortStringCount(senders, false)
	sortedRecips := sortStringCount(recips, true)
    result := SpamResults{sortedSubjs, sortedEds, sortedSenders, sortedRecips,
        msgCount}
    return result
}
