package email

import (
	"context"
	"net/smtp"

	log "github.com/jovanfrandika/smartbox-backend/pkg/common/logger"
)

type SendInput struct {
	To      string
	Subject string
	Body    string
}

const (
	emailHost = "smtp.gmail.com"
)

func (e *email) Send(ctx context.Context, sendInput SendInput) error {
	auth := smtp.PlainAuth("", e.config.EmailFrom, e.config.EmailPassword, emailHost)
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + sendInput.Subject + "!\n"
	msg := []byte(subject + mime + "\n" + sendInput.Body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, e.config.EmailFrom, []string{sendInput.To}, msg)
	if err != nil {
		log.Error(err.Error(), 0)
		return err
	}

	return nil
}
