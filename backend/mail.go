package main

import (
	"bytes"
	"context"
	"log"
	"time"

	gomail "github.com/wneessen/go-mail"
)

type Attachment struct {
	Name        string
	ContentType string
	Data        []byte
}

type Mail struct {
	To         string
	Subject    string
	Body       string
	Attachment *Attachment
}

// smtpTimeout limits one complete send.
const smtpTimeout = 30 * time.Second

type Mailer interface {
	Send(context.Context, Mail) error
}

// smtpMailer sends to a Microsoft 365 connector relay: STARTTLS, no login.
// The connector trusts the server's IP.
type smtpMailer struct {
	host string
	port int
	from string
}

func (m smtpMailer) Send(ctx context.Context, mail Mail) error {
	message, err := buildMessage(m.from, mail)
	if err != nil {
		return err
	}
	client, err := gomail.NewClient(m.host,
		gomail.WithPort(m.port),
		gomail.WithTLSPolicy(gomail.TLSMandatory),
		gomail.WithTimeout(smtpTimeout),
	)
	if err != nil {
		return err
	}
	// Finish the send even if the visitor leaves, but never take longer than smtpTimeout.
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), smtpTimeout)
	defer cancel()
	return client.DialAndSendWithContext(ctx, message)
}

// logMailer is used with DEV=1. It logs mails instead of sending them.
type logMailer struct{}

func (logMailer) Send(_ context.Context, m Mail) error {
	log.Printf("DEV mail to %s: %q (attachment: %v)\n%s", m.To, m.Subject, m.Attachment != nil, m.Body)
	return nil
}

func buildMessage(from string, mail Mail) (*gomail.Msg, error) {
	message := gomail.NewMsg()
	if err := message.From(from); err != nil {
		return nil, err
	}
	if err := message.To(mail.To); err != nil {
		return nil, err
	}
	message.Subject(mail.Subject)
	message.SetDate()
	message.SetBodyString(gomail.TypeTextPlain, mail.Body)
	if a := mail.Attachment; a != nil {
		err := message.AttachReader(a.Name, bytes.NewReader(a.Data), gomail.WithFileContentType(gomail.ContentType(a.ContentType)))
		if err != nil {
			return nil, err
		}
	}
	return message, nil
}
