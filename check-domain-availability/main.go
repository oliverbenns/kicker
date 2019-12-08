package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/domainr/whois"
	"log"
	"os"
	"strings"
)

func NotifyAvailable(domain string) {
	appName := os.Getenv("APP_NAME")
	message := domain + " is available."
	topicArn := os.Getenv("AWS_SNS_ARN")

	session := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion(os.Getenv("AWS_SNS_REGION"))
	svc := sns.New(session, config)

	svc.SetSMSAttributes(&sns.SetSMSAttributesInput{
		Attributes: map[string]*string{
			"SenderID": &appName,
		},
	})

	_, err := svc.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &topicArn,
	})

	if err != nil {
		log.Print("Error notifying", err)
	}
}

func Handler() {
	domains := strings.Split(os.Getenv("WANTED_DOMAINS"), ",")

	for _, domain := range domains {

		isAvailable := IsDomainAvailable(domain)

		if isAvailable {
			log.Print(domain, "is available.")
			NotifyAvailable(domain)
		}
	}
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
