package send

import (
	"fmt"
	"testing"
)

func Test_pairParticipants(t *testing.T) {
	type args struct {
		p Participants
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "Empty participants returns empty",
			args: args{
				p: map[string]Participant{},
			},
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name: "Two participants returns a pair",
			args: args{
				p: map[string]Participant{
					"1": {
						Name:      "1",
						Email:     "1@example.com",
						Interests: []string{"1"},
					},
					"2": {
						Name:      "2",
						Email:     "2@example.com",
						Interests: []string{"2"},
					},
				},
			},
			want: map[string]string{
				"1": "2",
				"2": "1",
			},
			wantErr: false,
		},
		{
			name: "unmatchable list returns an error",
			args: args{
				p: map[string]Participant{
					"1": {
						Name:      "1",
						Email:     "1@example.com",
						Interests: []string{"1"},
						Partner:   "2",
					},
					"2": {
						Name:      "2",
						Email:     "2@example.com",
						Interests: []string{"2"},
						Partner:   "1",
					},
				},
			},
			want:    map[string]string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pairParticipants(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("pairParticipants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for gifter, giftee := range got {
				if giftee.Name != tt.want[gifter.Name] {
					t.Errorf("pairParticipants() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

type testParticpantsLoaderError struct{}

func (t testParticpantsLoaderError) LoadParticipants(path string) (Participants, error) {
	return nil, fmt.Errorf("error loading participants")
}

type testParticpantsLoaderNoError struct {
	participantError bool
}

func (t testParticpantsLoaderNoError) LoadParticipants(path string) (Participants, error) {
	if t.participantError {
		// should return an error if the participants are unmatchable
		return map[string]Participant{
			"1": {
				Name:    "1",
				Email:   "",
				Partner: "2",
			},
			"2": {
				Name:    "2",
				Email:   "",
				Partner: "1",
			},
		}, nil
	}
	return map[string]Participant{
		"1": {
			Name:  "1",
			Email: "",
		},
		"2": {
			Name:  "2",
			Email: "",
		},
	}, nil
}

type testEmailerError struct{}

func (t testEmailerError) SendEmail(gifter, giftee Participant, emailTemplate *Email) error {
	return fmt.Errorf("error sending email")
}

type testEmailerNoError struct{}

func (t testEmailerNoError) SendEmail(gifter, giftee Participant, emailTemplate *Email) error {
	return nil
}

func TestSender_Send(t *testing.T) {
	type fields struct {
		Emailer           Emailer
		ParticipantLoader ParticipantLoader
		EmailTemplate     *Email
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Participant loader returns an error",
			fields: fields{
				Emailer:           nil,
				ParticipantLoader: testParticpantsLoaderError{},
				EmailTemplate:     &Email{},
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "Emailer returns an error",
			fields: fields{
				Emailer:           testEmailerError{},
				ParticipantLoader: testParticpantsLoaderNoError{},
				EmailTemplate:     &Email{},
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "PairParticipants returns an error",
			fields: fields{
				Emailer:           testEmailerNoError{},
				ParticipantLoader: testParticpantsLoaderNoError{participantError: true},
				EmailTemplate:     &Email{},
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "No errors",
			fields: fields{
				Emailer:           testEmailerNoError{},
				ParticipantLoader: testParticpantsLoaderNoError{participantError: false},
				EmailTemplate:     &Email{},
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sender{
				Emailer:           tt.fields.Emailer,
				ParticipantLoader: tt.fields.ParticipantLoader,
				EmailTemplate:     tt.fields.EmailTemplate,
			}
			if err := s.Send(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Sender.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
