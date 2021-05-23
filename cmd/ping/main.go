package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/oliverbenns/kicker/internal/notifications"
)

type Ctx struct {
	Notifier *notifications.Ctx
	GetUrls  func() []string
}

func (c *Ctx) Run() error {
	urls := c.GetUrls()
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("could not get website status: %w", err)
		}

		if resp.StatusCode >= 400 {
			message := url + " is down. Status code: " + strconv.Itoa(resp.StatusCode) + "."
			log.Print(message)
			err := c.Notifier.Notify(message)
			return err
		}
	}

	return nil
}

func Handler() {
	config := aws.NewConfig().WithRegion(os.Getenv("AWS_SNS_REGION"))
	session := session.Must(session.NewSession())

	notifierCtx := &notifications.Ctx{
		Sns:      sns.New(session, config),
		TopicArn: os.Getenv("AWS_SNS_ARN"),
	}

	ctx := Ctx{
		Notifier: notifierCtx,
		GetUrls: func() []string {
			return strings.Split(os.Getenv("CONCERNED_ENDPOINTS"), ",")

		},
	}

	ctx.Run()
}

func main() {
	lambda.Start(Handler)
}
