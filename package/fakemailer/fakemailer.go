package fakemailer

import (
	"fmt"

	"github.com/dcmcand/go-secret-santa/package/send"
)

// A fake mailer for testing purposes. It does not actually send emails.
// it just logs the email that would have been sent.

type Mailer struct{}

func (m *Mailer) SendEmail(gifter, giftee send.Participant, emailTemplate *send.Email) error {
	mail, err := emailTemplate.Render(gifter, giftee)
	if err != nil {
		return fmt.Errorf("error rendering email: %v", err)
	}
	fmt.Printf("\nEmail to %s <%s>:\n%s\n", gifter.Name, gifter.Email, mail)
	return nil
}
