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

func NewSns() (*testSns, error) {
	config := &aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String("http://localhost:4566"),
	}

	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	snsClient := sns.New(session, config)

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
