package email_service

import (
	"github.com/alonelegion/go_graphql_api/email_client/mailgun_client"
)

type EmailService interface {
	Welcome(toEmail string) error
	ResetPassword(toEmail, token string) error
}

type emailService struct {
	client mailgun_client.MailGunClient
}

func NewEmailService(client mailgun_client.MailGunClient) EmailService {
	return &emailService{
		client: client,
	}
}

func (e emailService) Welcome(toEmail string) error {
	return e.client.Welcome(welcomeTheme, welcomeText, toEmail, welcomeHTML)
}

func (e *emailService) ResetPassword(toEmail, token string) error {
	return e.client.ResetPassword(resetTheme, resetTextTmpl, toEmail, resetHTMLTmpl, token)
}
