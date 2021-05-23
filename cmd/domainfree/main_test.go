package main

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/oliverbenns/kicker/internal/notifications"
	"github.com/oliverbenns/kicker/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testBucketName = "test-bucket"

func createTestBucket(sess *session.Session) error {
	s3Client := s3.New(sess)

	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(testBucketName),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeBucketAlreadyExists {
			return nil
		}
	}

	return err
}

func uploadToBucket(sess *session.Session) error {
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	buf := bytes.NewBufferString("name,com,co\ntzacwierjiyknoelkefbmyankdnlxbvaoujuizfy,1,0")

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String("domains.csv"),
		Body:   buf,
	})

	return err
}

func TestDomainFreeNotifies(t *testing.T) {
	err := test.WaitForLocalStack()
	require.NoError(t, err)

	sess, err := test.NewAwsSession()
	require.NoError(t, err)

	err = createTestBucket(sess)
	require.NoError(t, err)

	err = uploadToBucket(sess)
	require.NoError(t, err)

	downloader := s3manager.NewDownloader(sess)

	sns, err := test.NewSns(sess)
	require.NoError(t, err)

	sub, err := test.NewSubscriber(sns)
	require.NoError(t, err)

	go func() {
		err := sub.Listen()
		require.NoError(t, err)
	}()

	notifierCtx := &notifications.Ctx{
		Sns:      sns.Client,
		TopicArn: *sns.TopicArn,
	}

	ctx := Ctx{
		Notifier:   notifierCtx,
		Downloader: downloader,
		BucketName: testBucketName,
	}

	ctx.Run()

	notification, err := sub.GetNotification()
	require.NoError(t, err)
	assert.Equal(t, "tzacwierjiyknoelkefbmyankdnlxbvaoujuizfy.com is available.", notification.Message)
}
