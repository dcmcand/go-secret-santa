package mgmailer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dcmcand/go-secret-santa/package/send"
	"github.com/mailgun/mailgun-go/v4"
)

type MailgunEmailer struct {
	emailTemplate *send.Email
	mg            mailgun.Mailgun
}

func (m *MailgunEmailer) SendEmail(gifter, giftee send.Participant, emailTemplate *send.Email) error {
	// The message object allows you to add attachments and Bcc recipients
	body, err := emailTemplate.Render(gifter, giftee)
	if err != nil {
		return fmt.Errorf("error rendering email: %v", err)
	}
	message := mailgun.NewMessage(emailTemplate.SenderEmail, emailTemplate.Subject, body, gifter.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := m.mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return nil
}

func NewMailgunEmailer(domain, apiKey string) *MailgunEmailer {
	return &MailgunEmailer{
		mg: mailgun.NewMailgun(domain, apiKey),
	}
}
