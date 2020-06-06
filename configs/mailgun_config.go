package configs

import "os"

type MailGunConfig struct {
	APIKey string `env:"MAILGUN_API_KEY"`

	Domain string `env:"MAILGUN_DOMAIN"`
}

func GetMailGunConfig() MailGunConfig {
	return MailGunConfig{
		APIKey: os.Getenv("MAILGUN_API_KEY"),

		Domain: os.Getenv("MAILGUN_DOMAIN"),
	}
}
