package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oliverbenns/kicker/internal/notifications"
)

type Ctx struct {
	Notify  func(domain string)
	GetUrls func() []string
}

func (c *Ctx) Run() {
	urls := c.GetUrls()
	for _, url := range urls {
		resp, err := http.Get(url)

		if err != nil {
			log.Print("Error getting website status", url, err)
		}

		if resp.StatusCode >= 400 {
			message := url + " is down. Status code: " + strconv.Itoa(resp.StatusCode) + "."
			log.Print(message)
			c.Notify(message)
		}
	}
}

func Handler() {
	ctx := Ctx{
		Notify: notifications.Notify,
		GetUrls: func() []string {
			return strings.Split(os.Getenv("CONCERNED_ENDPOINTS"), ",")

		},
	}

	ctx.Run()
}

func main() {
	lambda.Start(Handler)
}
