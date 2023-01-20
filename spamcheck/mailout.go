package main

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	acctname = "ed"
	from     = "ed@leafe.com"
	hostname = "localhost"
	hostport = "localhost:2525"
	msg      = []byte("dummy message")
	pw       = os.Getenv("MAILPW")
	//	recipients = []string{"ed@leafe.com", "edleafe@gmail.com"}
	recipients = []string{"edleafe@gmail.com"}
)

func makeLinkText(cnt int, loc string) string {
	link := "https://mailmaint.leafe.com/spam/" + filepath.Base(loc)
	return fmt.Sprintf(`Filtered Message Total: %d
To Delete: <%s>`, cnt, link)
}

func makeHeader(results SpamResults, loc, prefix string) string {
	tm := time.Now()
	tmStr := tm.Format("3:04 PM on Jan 2, 2006")
	prf := ""
	if prefix != "" {
		prf = fmt.Sprintf("%s ", strings.Title(prefix))
	}
	template := `To: Ed Leafe <ed@leafe.com>
Subject: %sSpam Check - Go

%sSpam Header Check
Time: %s
Spam File: %s`
	msgHeader := fmt.Sprintf(template, prf, prf, tmStr, loc)
	return msgHeader
}

func joinSlice(slc []string) string {
	return strings.Join(slc, "\n")
}

func assemble(results SpamResults, loc, prefix string) string {
	linkText := makeLinkText(results.Count, loc)
	header := makeHeader(results, loc, prefix)

	text := `%s

%s

Recipients:
============
%s

%s

Subjects:
==========
%s

%s

ED Subjects:
============
%s

Senders:
==========
%s

%s
`
	return fmt.Sprintf(text, header, linkText, joinSlice(results.SortedRecips),
		linkText, joinSlice(results.SortedSubjects), linkText,
		joinSlice(results.SortedEds), joinSlice(results.SortedSenders),
		linkText)
}

func MailOut(results SpamResults, loc, prefix string) {
	content := assemble(results, loc, prefix)

	c, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	// Set the sender and recipient.
	c.Mail("spamCheck@leafe.com")
	c.Rcpt("ed@leafe.com")
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString(content)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}
