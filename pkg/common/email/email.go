package email

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	log "github.com/jovanfrandika/smartbox-backend/pkg/common/logger"
)

type SendInput struct {
	To      string
	Subject string
	Body    string
}

func (e *email) Send(ctx context.Context, sendInput SendInput) error {
	_, err := e.client.SendEmail(ctx, &ses.SendEmailInput{
		Source: &e.config.EmailIdentitySource,
		Destination: &types.Destination{
			ToAddresses: []string{sendInput.To},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: &sendInput.Subject,
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: &sendInput.Body,
				},
			},
		},
	})
	if err != nil {
		log.Error(err.Error(), 0)
		return err
	}

	return nil
}
