package email_service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mgClient struct {
	mock.Mock
}

func (cl *mgClient) Welcome(subject, text, to, htmlStr string) error {
	args := cl.Called(subject, text, to, htmlStr)
	return args.Error(0)
}

func (cl *mgClient) ResetPassword(subject, text, to, htmlStr, token string) error {
	args := cl.Called(subject, text, to, htmlStr, token)
	return args.Error(0)
}

func TestEmailService(t *testing.T) {
	client := new(mgClient)
	es := NewEmailService(client)

	t.Run("Welcome", func(t *testing.T) {
		toEmail := "example@mail.org"
		client.On("Welcome", welcomeTheme, welcomeText, toEmail, welcomeHTML).Return(nil)

		err := es.Welcome(toEmail)
		assert.Nil(t, err)
	})

	t.Run("ResetPassword", func(t *testing.T) {
		toEmail := "example@mail.org"
		token := "secret-token"
		client.On("ResetPassword", resetTheme, resetTextTmpl, toEmail, resetHTMLTmpl, token).Return(nil)

		err := es.ResetPassword(toEmail, token)
		assert.Nil(t, err)
	})
}
