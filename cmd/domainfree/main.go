package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/domainr/whois"
	"github.com/oliverbenns/kicker/internal/notifications"
)

type Ctx struct {
	Notifier   *notifications.Ctx
	Downloader *s3manager.Downloader
	BucketName string
}

func (c *Ctx) DownloadCsv() (*bytes.Buffer, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := c.Downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(c.BucketName),
		Key:    aws.String("domains.csv"),
	})

	return bytes.NewBuffer(buf.Bytes()), err
}

func (c *Ctx) Run() error {
	data, err := c.DownloadCsv()
	if err != nil {
		return fmt.Errorf("error obtaining domain list: %w", err)
	}

	r := csv.NewReader(data)
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
	sess := session.Must(session.NewSession())

	downloader := s3manager.NewDownloader(sess)

	notifierCtx := &notifications.Ctx{
		Sns:      sns.New(sess),
		TopicArn: os.Getenv("AWS_SNS_ARN"),
	}

	ctx := Ctx{
		Notifier:   notifierCtx,
		Downloader: downloader,
		BucketName: "kicker-data",
	}

	ctx.Run()
}

func main() {
	lambda.Start(Handler)
}
