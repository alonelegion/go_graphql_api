package mailgun_client

import (
	"context"
	"fmt"
	"github.com/alonelegion/go_graphql_api/configs"
	"github.com/mailgun/mailgun-go/v4"
	"net/url"
	"time"
)

type MailGunClient interface {
	Welcome(theme, text, to, htmlStr string) error
	ResetPassword(theme, text, to, htmlStr, token string) error
}

type mailGunClient struct {
	config configs.Config
	client *mailgun.MailgunImpl
}

func NewMailGunClient(c configs.Config) MailGunClient {
	return &mailGunClient{
		config: c,
		client: mailgun.NewMailgun(c.MailGun.Domain, c.MailGun.APIKey),
	}
}

func (m *mailGunClient) Welcome(theme, text, to, htmlStr string) error {
	message := m.createNewMessage(
		m.config.FromEmail,
		theme,
		text,
		to,
		htmlStr,
	)

	ctx, cancel := m.setContext(10)
	defer cancel()
	return m.send(ctx, message)
}

func (m *mailGunClient) ResetPassword(theme, text, to, htmlStr, token string) error {
	v := url.Values{}
	v.Set("token", token)

	resetURL := m.getURL() + "/apiserver/update_password?" + v.Encode()
	resetText := fmt.Sprintf(text, resetURL, token)
	resetHTML := fmt.Sprintf(htmlStr, resetURL, token)
	message := m.createNewMessage(
		m.config.FromEmail,
		theme,
		resetText,
		to,
		resetHTML,
	)

	ctx, cancel := m.setContext(10)
	defer cancel()
	return m.send(ctx, message)
}

func (m *mailGunClient) getURL() string {
	url := m.config.Host + ":" + m.config.Port
	return url
}

func (m *mailGunClient) createNewMessage(from, theme, text, to, htmlStr string) *mailgun.Message {
	message := m.client.NewMessage(
		from,
		theme,
		text,
		to,
	)
	message.SetHtml(htmlStr)
	return message
}

func (m *mailGunClient) setContext(second time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*second)
}

func (m *mailGunClient) send(ctx context.Context, message *mailgun.Message) error {
	_, _, err := m.client.Send(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
