package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type testSns struct {
	Client   *sns.SNS
	TopicArn *string
}

func NewSns(sess *session.Session) (*testSns, error) {
	snsClient := sns.New(sess)

	topic, err := snsClient.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String("Kicker"),
	})
	if err != nil {
		return nil, err
	}

	return &testSns{
		Client:   snsClient,
		TopicArn: topic.TopicArn,
	}, nil
}
