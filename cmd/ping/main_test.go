package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/oliverbenns/kicker/internal/notifications"
	"github.com/oliverbenns/kicker/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const data = "url\nhttps://google.com/404"
const fileName = "urls.csv"

func TestPingNotifies(t *testing.T) {
	err := test.WaitForLocalStack()
	require.NoError(t, err)

	sess, err := test.NewAwsSession()
	require.NoError(t, err)

	bucketName, err := test.CreateBucket(sess)
	require.NoError(t, err)

	err = test.UploadToBucket(sess, fileName, data)
	require.NoError(t, err)

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
		Downloader: s3manager.NewDownloader(sess),
		BucketName: bucketName,
	}

	ctx.Run()

	notification, err := sub.GetNotification()
	require.NoError(t, err)
	assert.Equal(t, "https://google.com/404 is down. Status code: 404.", notification.Message)
}
