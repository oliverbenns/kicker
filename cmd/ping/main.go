package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
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
		Key:    aws.String("urls.csv"),
	})

	return bytes.NewBuffer(buf.Bytes()), err
}

func (c *Ctx) Run() error {
	data, err := c.DownloadCsv()
	if err != nil {
		return fmt.Errorf("error obtaining url list: %w", err)
	}

	r := csv.NewReader(data)

	// headers
	_, err = r.Read()
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

		if len(row) != 1 {
			return errors.New("bad amount of vat rows")
		}

		url := row[0]

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("could not get website status: %w", err)
		}

		if resp.StatusCode >= 400 {
			message := url + " is down. Status code: " + strconv.Itoa(resp.StatusCode) + "."
			log.Print(message)
			return c.Notifier.Notify(message)
		}
	}

	return nil
}

func Handler() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	notifierCtx := &notifications.Ctx{
		Sns:      sns.New(sess),
		TopicArn: os.Getenv("AWS_SNS_ARN"),
	}

	ctx := Ctx{
		Notifier:   notifierCtx,
		Downloader: s3manager.NewDownloader(sess),
		BucketName: "kicker-data",
	}

	return ctx.Run()
}

func main() {
	lambda.Start(Handler)
}
