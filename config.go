package main

// Config defines configuration variables with their mapping from environment variables.
// It also has default values for running dev application.
type Config struct {
	Listen                   string `env:"LISTEN" envDefault:":8080"`
	HTTPProxy                string `env:"HTTP_PROXY" envDefault:""`
	DoerTimeoutSeconds       int    `env:"DOER_TIMEOUT_SECONDS" envDefault:"10"`
	HTTPServerTimeoutSeconds int    `env:"SERVER_TIMEOUT_SECONDS" envDefault:"30"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`

	MetricsHost    string `env:"METRICS_HOST" envDefault:"localhost"`
	MetricsPort    string `env:"METRICS_PORT" envDefault:"8125"`
	MetricsEnabled bool   `env:"METRICS_ENABLED" envDefault:"false"`
	MetricsPrefix  string `env:"METRICS_PREFIX" envDefault:"server-api"`

	MailGunHost       string `env:"MAIL_MAILGUN_HOST" envDefault:"api.mailgun.net"`
	MailGunAPIKey     string `env:"MAIL_MAILGUN_API_KEY" envDefault:""`
	MailGunDomainName string `env:"MAIL_MAILGUN_DOMAIN_NAME" envDefault:""`

	SendGridHost    string `env:"MAIL_SENDGRID_HOST" envDefault:"api.sendgrid.com"`
	SendGridAPIKey  string `env:"MAIL_SENDGRID_API_KEY" envDefault:""`
	SendGridAPIUser string `env:"MAIL_SENDGRID_API_USER" envDefault:""`
}
