package email

import (
	"log/slog"

	"github.com/resend/resend-go/v3"
)

type EmailConfig struct {
	apiKey string
	logger *slog.Logger
}

func NewEmailConfig(apiKey string, logger *slog.Logger) *EmailConfig {
	return &EmailConfig{
		apiKey: apiKey,
		logger: logger,
	}
}

func (e *EmailConfig) SendEmail(to, body string) error {
	client := resend.NewClient(e.apiKey)

	params := &resend.SendEmailRequest{
		From:    "Anurag Singh <send@anuragcode.me>",
		Subject: "From Carbon",
		To:      []string{to},
		Html:    body,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		e.logger.Error("error while sending email", "error", err)
		return err
	}

	return nil
}
