package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/oliverbenns/kicker/internal/notifications"
	"github.com/oliverbenns/kicker/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainFreeNotifies(t *testing.T) {
	err := test.WaitForLocalStack()
	require.NoError(t, err)

	go func() {
		http.HandleFunc("/domains.csv", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "name,com,co\ntzacwierjiyknoelkefbmyankdnlxbvaoujuizfy,1,0")
		})
		http.ListenAndServe(":8080", nil)
	}()

	sns, err := test.NewSns()
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
		GetCsvUrl: func() string {
			return "http://localhost:8080/domains.csv"
		},
	}

	ctx.Run()

	notification, err := sub.GetNotification()
	require.NoError(t, err)
	assert.Equal(t, "tzacwierjiyknoelkefbmyankdnlxbvaoujuizfy.com is available.", notification.Message)
}
