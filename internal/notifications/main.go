package notifications

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sns"
)

type Ctx struct {
	Sns      *sns.SNS
	TopicArn string
}

func (c *Ctx) Notify(message string) error {
	_, err := c.Sns.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &c.TopicArn,
	})

	if err != nil {
		return fmt.Errorf("failed to notify: %w", err)
	}

	return nil
}
