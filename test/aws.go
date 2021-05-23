package test

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func NewAwsSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:           aws.String("eu-west-1"),
		Endpoint:         aws.String("http://localhost:4566"),
		S3ForcePathStyle: aws.Bool(true),
	})
}

const bucketName = "test-bucket"

func CreateBucket(sess *session.Session) (string, error) {
	s3Client := s3.New(sess)

	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeBucketAlreadyExists {
			return bucketName, nil
		}
	}

	return bucketName, err
}

func UploadToBucket(sess *session.Session, fileName, data string) error {
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	buf := bytes.NewBufferString(data)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   buf,
	})

	return err
}
