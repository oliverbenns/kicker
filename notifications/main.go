package notifications

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
	"os"
)

type Notifier = func(domain string)

func Notify(message string) {
	appName := os.Getenv("APP_NAME")
	topicArn := os.Getenv("AWS_SNS_ARN")

	session := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion(os.Getenv("AWS_SNS_REGION"))
	svc := sns.New(session, config)

	svc.SetSMSAttributes(&sns.SetSMSAttributesInput{
		Attributes: map[string]*string{
			"SenderID": &appName,
		},
	})

	_, err := svc.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &topicArn,
	})

	if err != nil {
		log.Print("Error notifying", err)
	}
}
