package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oliverbenns/kicker/notifications"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func CreateHandler(notify notifications.Notifier) func() {
	return func() {
		endpoints := strings.Split(os.Getenv("CONCERNED_ENDPOINTS"), ",")

		for _, endpoint := range endpoints {
			resp, err := http.Get(endpoint)

			if err != nil {
				log.Print("Error getting website status", endpoint, err)
			}

			if resp.StatusCode >= 400 {
				message := endpoint + " is down. Status code: " + strconv.Itoa(resp.StatusCode) + "."
				log.Print(message)
				notify(message)
			}
		}
	}
}

// @TODO: Is there a better way to do this? func Handler() = CreateHandler(NotifyAvailable)
func Handler() {
	handler := CreateHandler(notifications.Notify)
	handler()
}

func main() {
	lambda.Start(Handler)
}
