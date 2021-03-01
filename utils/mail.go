package utils

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
)

// GMail smtp server
const (
	SMTPServer = "smtp.pobox.com"
)

// Authorizer is the sender of the message
type Authorizer struct {
	User     string
	Password string
}

// NewAuthorizer creates a new sender
func NewAuthorizer(Username, Password string) Authorizer {
	return Authorizer{Username, Password}
}

// SendMail sends a basic message
func (a Authorizer) SendMail(Dest []string, sender, Subject, bodyMessage string) {

	msg := "From: " + sender + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err := smtp.SendMail(SMTPServer+":587",
		smtp.PlainAuth("", a.User, a.Password, SMTPServer),
		sender, Dest, []byte(msg))

	if err != nil {

		fmt.Printf("smtp error: %s", err)
		return
	}

	fmt.Println("Mail sent successfully!")
}

func (a Authorizer) writeEmail(dest []string, sender, contentType, subject, bodyMessage string) string {

	header := make(map[string]string)
	header["From"] = sender

	receipient := ""

	for _, user := range dest {
		receipient = receipient + user
	}

	header["To"] = receipient
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", contentType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer

	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	finalMessage.Write([]byte(bodyMessage))
	finalMessage.Close()

	message += "\r\n" + encodedMessage.String()

	return message
}

// WriteHTMLEmail composes the body of the message as an HTML email
func (a *Authorizer) WriteHTMLEmail(dest []string, from, subject, bodyMessage string) string {
	return a.writeEmail(dest, from, "text/html", subject, bodyMessage)
}

// WritePlainEmail composes the body of the message as plain text
func (a *Authorizer) WritePlainEmail(dest []string, from, subject, bodyMessage string) string {
	return a.writeEmail(dest, from, "text/plain", subject, bodyMessage)
}
