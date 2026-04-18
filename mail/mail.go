package mail

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
)

// Message is a simple email.
type Message struct {
	From    string
	To      []string
	Subject string
	Body    string
	HTML    bool
}

// SMTP sends via smtp.SendMail.
type SMTP struct {
	Addr string
	Auth smtp.Auth
}

// Send delivers message.
func (s *SMTP) Send(ctx context.Context, m *Message) error {
	_ = ctx
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: %s\r\n", m.From)
	if len(m.To) == 0 {
		return fmt.Errorf("mail: no recipients")
	}
	fmt.Fprintf(&buf, "To: %s\r\n", m.To[0])
	fmt.Fprintf(&buf, "Subject: %s\r\n", m.Subject)
	if m.HTML {
		fmt.Fprintf(&buf, "Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		fmt.Fprintf(&buf, "Content-Type: text/plain; charset=UTF-8\r\n")
	}
	fmt.Fprintf(&buf, "\r\n%s", m.Body)
	return smtp.SendMail(s.Addr, s.Auth, m.From, m.To, buf.Bytes())
}
