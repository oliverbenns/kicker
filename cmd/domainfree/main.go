package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/domainr/whois"
	"github.com/oliverbenns/kicker/internal/notifications"
)

type Ctx struct {
	Notifier  *notifications.Ctx
	GetCsvUrl func() string
}

func (c *Ctx) Run() error {
	url := c.GetCsvUrl()
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error obtaining url list: %w", err)
	}

	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)

	headers, err := r.Read()
	if err != nil {
		return fmt.Errorf("error reading csv: %w", err)
	}

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading csv row: %w", err)
		}

		var name string
		for i, column := range row {
			if i == 0 {
				name = row[i]
				continue
			}
			val, err := strconv.Atoi(column)

			if err != nil {
				return fmt.Errorf("error converting column: %w", err)
			}

			if val > 0 {
				tld := headers[i]
				domain := name + "." + tld

				if IsDomainAvailable(domain) {
					message := domain + " is available."
					log.Print(message)
					err := c.Notifier.Notify(message)
					return err
				}
			}
		}
	}

	return nil
}

func IsDomainAvailable(domain string) bool {
	request, _ := whois.NewRequest(domain)

	out, err := whois.DefaultClient.Fetch(request)

	if err != nil {
		log.Print("Error performing whois query", domain, err)
		return false
	}

	// No match for - most domains
	// No Data Found - .co
	str := string(out.Body)
	return strings.Contains(str, "No match for") || strings.Contains(str, "No Data Found")
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
		GetCsvUrl: func() string {
			return os.Getenv("DOMAINS_CSV_URL")
		},
	}

	err := ctx.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	lambda.Start(Handler)
}
