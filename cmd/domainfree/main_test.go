package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/oliverbenns/kicker/internal/notifications"
	"github.com/oliverbenns/kicker/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const data = "name,com,co\ntzacwierjiyknoelkefbmyankdnlxbvaoujuizfy,1,0"
const fileName = "domains.csv"

func TestDomainFreeNotifies(t *testing.T) {
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
	assert.Equal(t, "tzacwierjiyknoelkefbmyankdnlxbvaoujuizfy.com is available.", notification.Message)
}
