package email

type EmailConfig struct {
	ApiKey string
}

func (e *EmailConfig) SendEmail(to, body string) {
}
