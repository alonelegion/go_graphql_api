package configs

import "os"

const (
	production = "production"
)

type Config struct {
	Env       string         `env:"ENV"`
	Pepper    string         `env:"PEPPER"`
	HMACKey   string         `env:"HMAC_KEY"`
	Postgres  PostgresConfig `json:"postgres"`
	MailGun   MailGunConfig  `json:"mailgun"`
	JWTSecret string         `env:"JWT_SIGN_KEY"`
	Host      string         `env:"APP_HOST"`
	Port      string         `env:"APP_PORT"`
	FromEmail string         `env:"EMAIL_FROM"`
}

func (c Config) IsProduction() bool {
	return c.Env == production
}

func GetConfig() Config {
	return Config{
		Env:       os.Getenv("ENV"),
		Pepper:    os.Getenv("PEPPER"),
		HMACKey:   os.Getenv("HMAC_KEY"),
		Postgres:  GetPostgresConfig(),
		MailGun:   GetMailGunConfig(),
		JWTSecret: os.Getenv("JWT_SIGN_KEY"),
		Host:      os.Getenv("APP_HOST"),
		Port:      os.Getenv("APP_PORT"),
		FromEmail: os.Getenv("EMAIL_FROM"),
	}
}
