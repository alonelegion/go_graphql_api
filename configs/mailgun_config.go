package configs

import "os"

type MailgunConfig struct {
	APIKey string `env:"MAILGUN_API_KEY"`

	Domain string `env:"MAILGUN_DOMAIN"`
}

func GetMailgunConfig() MailgunConfig {
	return MailgunConfig{
		APIKey: os.Getenv("MAILGUN_API_KEY"),

		Domain: os.Getenv("MAILGUN_DOMAIN"),
	}
}
