package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewAwsSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:           aws.String("eu-west-1"),
		Endpoint:         aws.String("http://localhost:4566"),
		S3ForcePathStyle: aws.Bool(true),
	})
}
