package main

import (
	"testing"

	"github.com/oliverbenns/kicker/internal/notifications"
	"github.com/oliverbenns/kicker/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingNotifies(t *testing.T) {
	err := test.WaitForLocalStack()
	require.NoError(t, err)

	sess, err := test.NewAwsSession()
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
		Notifier: notifierCtx,
		GetUrls: func() []string {
			return []string{"https://google.com/404"}
		},
	}

	ctx.Run()

	notification, err := sub.GetNotification()
	require.NoError(t, err)
	assert.Equal(t, "https://google.com/404 is down. Status code: 404.", notification.Message)
}
