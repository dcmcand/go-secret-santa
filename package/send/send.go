package send

import (
	"bytes"
	"fmt"
	"text/template"
)

type Emailer interface {
	SendEmail(gifter, giftee Participant, emailTemplate *Email) error
}

type ParticipantLoader interface {
	LoadParticipants(path string) (Participants, error)
}

type Email struct {
	Subject     string
	Body        *template.Template
	SenderName  string
	SenderEmail string
}

func (e Email) Render(gifter, giftee Participant) (string, error) {
	var buff bytes.Buffer
	err := e.Body.Execute(&buff, struct {
		Gifter Participant
		Giftee Participant
	}{
		Gifter: gifter,
		Giftee: giftee,
	})
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}
	return buff.String(), nil
}

type Sender struct {
	Emailer           Emailer
	ParticipantLoader ParticipantLoader
	EmailTemplate     *Email
}

type Participant struct {
	Name      string
	Email     string
	Interests []string
	Partner   string
}

type Participants map[string]Participant

func (s *Sender) Send(path string) error {
	participants, err := s.ParticipantLoader.LoadParticipants(path)
	if err != nil {
		return fmt.Errorf("error parsing participants: %v", err)
	}
	pairs, err := pairParticipants(participants)
	if err != nil {
		return fmt.Errorf("error pairing participants: %v", err)
	}
	for gifter, giftee := range pairs {
		err = s.Emailer.SendEmail(*gifter, *giftee, s.EmailTemplate)
	}
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}

type pairedParticipants map[*Participant]*Participant

func pairParticipants(p Participants) (pairedParticipants, error) {

	names := make(map[string]struct{})
	for name, _ := range p {
		names[name] = struct{}{}
	}
	pairs := make(pairedParticipants)
	// Since this is non-deterministic, we'll try a few times
	var err error
	for tries := 10; tries >= 0; tries-- {
		for _, participant := range p {
			var gifteeName string
			gifteeName, err = getName(participant, names)
			// don't continue if we can't find a name
			if err != nil {
				break
			}
			giftee := p[gifteeName]
			pairs[&participant] = &giftee
			delete(names, gifteeName)
		}
		// don't retry if we found a solution
		if err == nil {
			break
		}
		// return error if tries are exhausted
		if tries == 0 {
			return pairedParticipants{}, fmt.Errorf("error finding giftee: %v", err)
		}
		// reset names
		for name, _ := range p {
			names[name] = struct{}{}
		}

	}
	return pairs, nil
}

func getName(gifter Participant, names map[string]struct{}) (string, error) {

	name := ""
	for key, _ := range names {
		if key != gifter.Name && key != gifter.Partner {
			name = key
			break
		}
	}
	if name == "" {
		return name, fmt.Errorf("no name found")
	}
	return name, nil
}
