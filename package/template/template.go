package template

import (
	"fmt"
	"text/template"

	"github.com/dcmcand/go-secret-santa/package/send"
)

func GetDefaultTemplate(subject, senderName, senderEmail string) (*send.Email, error) {
	fmt.Println("getting default template")
	tmplSrc := `Hello {{.Gifter.Name}},
This is your secret santa assignment!
This Christmas, you will buy a gift for {{.Giftee.Name}}.
{{.Giftee.Name}} wrote in their letter to Santa that they are interested in {{.Giftee.Interests}}.
Remember this is a SECRET Santa so ssssshhhhhhh!
Merry Christmas
Santa Claus`
	tmpl, err := template.New("default").Parse(tmplSrc)
	if err != nil {
		return nil, err
	}
	return &send.Email{
		Subject:     subject,
		SenderName:  senderName,
		SenderEmail: senderEmail,
		Body:        tmpl,
	}, nil
}

func GetTemplate(tmplSrc, subject, senderName, senderEmail string) (*send.Email, error) {
	tmpl, err := template.New("custom").ParseFiles(tmplSrc)
	if err != nil {
		return nil, err
	}
	return &send.Email{
		Subject:     subject,
		SenderName:  senderName,
		SenderEmail: senderEmail,
		Body:        tmpl,
	}, nil

}
