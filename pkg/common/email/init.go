package email

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
)

type email struct {
	client *ses.Client
	config *config.Config
}

type Email interface {
	Send(ctx context.Context, sendInput SendInput) error
}

func New(awsConfig *aws.Config, config *config.Config) Email {
	client := ses.NewFromConfig(*awsConfig)
	return &email{
		client: client,
		config: config,
	}
}
