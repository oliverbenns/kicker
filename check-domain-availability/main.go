package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/domainr/whois"
	"log"
	"os"
	"strings"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	var domainsAvailable string
	var domainsUnavailable string
	domains := strings.Split(os.Getenv("WANTED_DOMAINS"), ",")

	for _, domain := range domains {

		isAvailable := IsDomainAvailable(domain)

		if isAvailable {
			domainsAvailable += domain
		} else {
			domainsUnavailable += domain
		}
	}

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message":               "Go Serverless v1.0! Your function executed successfully!",
		"domains available":     domainsAvailable,
		"domains not available": domainsUnavailable,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "check-domain-availability-handler",
		},
	}

	return resp, nil
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
