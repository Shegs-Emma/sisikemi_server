package mail

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailSendGridSender interface {
	SendGridEmail(subject, content, toEmail string) error
}

type SendGridEmailSender struct {
	apiKey string
	fromEmail string
}

func NewSendGridEmailSender (apiKey, fromEmail string) *SendGridEmailSender {
	return &SendGridEmailSender{
		apiKey: apiKey,
		fromEmail: fromEmail,
	}
}

func (s *SendGridEmailSender) SendGridEmail(subject, content, toEmail string) error {
	from := mail.NewEmail("Sisikemi Fashion", s.fromEmail)

	fmt.Println("from", from)
	to := mail.NewEmail("", toEmail)
	
	fmt.Println("to", to)
	message := mail.NewSingleEmail(from, subject, to, content, content)

	fmt.Println("message", message)
	client := sendgrid.NewSendClient(s.apiKey)

	fmt.Println("client", client)
	_, err := client.Send(message)

	fmt.Println("err", err)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}