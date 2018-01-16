package main

import (
	"fmt"
	"log"
    "net/smtp"
	"os"
    "path/filepath"
    "strings"
	"time"
)

var (
    acctname   = "ed"
	from       = "ed@leafe.com"
	hostname   = "mail.leafe.com"
	hostport   = "mail.leafe.com:2525"
	msg        = []byte("dummy message")
	pw         = os.Getenv("MAILPW")
//	recipients = []string{"ed@leafe.com", "edleafe@gmail.com"}
	recipients = []string{"edleafe@gmail.com"}
)

func makeLinkText(cnt int, loc string) string {
    link := "http://mail.leafe.com/cgi-bin/delspam/" + filepath.Base(loc)
    return fmt.Sprintf(`Filtered Message Total: %d
To Delete: <%s>`, cnt, link)
}

func makeHeader(results SpamResults, loc string) string {
    tm := time.Now()
    tmStr := tm.Format("3:04 PM on Jan 2, 2006")

    // The first format field will eventually be varible, but for now, just
    // leave it blank.
    msgHeader := fmt.Sprintf(`%sSpam Header Check
Time: %s
Spam File: %s`, "", tmStr, loc)
    return msgHeader
}

func joinSlice(slc []string) string {
    return strings.Join(slc, "\n")
}

func assemble(results SpamResults, loc string) string {
    linkText := makeLinkText(results.Count, loc)
    header := makeHeader(results, loc)

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

func MailOut(results SpamResults, loc string) {
    content := assemble(results, loc)

	// Set up authentication information.
	auth := smtp.PlainAuth("", acctname, pw, hostname)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	msg := []byte("To: edleafe@gmail.com\r\n" +
		"From: ed@leafe.com\r\n" +
		"Subject: SpamTest results\r\n" +
		"\r\n" + content + "\r\n")
	err := smtp.SendMail(hostport, auth, from, recipients, msg)
	if err != nil {
		log.Fatal(err)
	}
}
