package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/domainr/whois"
	"github.com/oliverbenns/kicker/notifications"
	"log"
	"os"
	"strings"
)

func CreateHandler(notify notifications.Notifier) func() {
	return func() {
		domains := strings.Split(os.Getenv("WANTED_DOMAINS"), ",")

		for _, domain := range domains {
			if IsDomainAvailable(domain) {
				message := domain + " is available."
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

func IsDomainAvailable(domain string) bool {
	request, _ := whois.NewRequest(domain)

	out, err := whois.DefaultClient.Fetch(request)

	if err != nil {
		log.Print("Error performing whois query", domain, err)
		return false
	}

	return strings.Contains(string(out.Body), "No match for")
}

func main() {
	lambda.Start(Handler)
}
