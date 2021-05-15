package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/domainr/whois"
	"github.com/oliverbenns/kicker/internal/notifications"
)

type Ctx struct {
	Notify    func(domain string)
	GetCsvUrl func() string
}

func (c *Ctx) Run() {
	url := c.GetCsvUrl()
	resp, err := http.Get(url)

	if err != nil {
		panic("Error obtaining url list")
	}

	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)

	headers, err := r.Read()

	if err != nil {
		log.Fatal(err)
	}

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		var name string
		for i, column := range row {
			if i == 0 {
				name = row[i]
				continue
			}
			val, err := strconv.Atoi(column)

			if err != nil {
				log.Fatal(err)
			}

			if val > 0 {
				tld := headers[i]
				domain := name + "." + tld

				if IsDomainAvailable(domain) {
					message := domain + " is available."
					log.Print(message)
					c.Notify(message)
				}
			}
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

	// No match for - most domains
	// No Data Found - .co
	str := string(out.Body)
	return strings.Contains(str, "No match for") || strings.Contains(str, "No Data Found")
}

func Handler() {
	ctx := Ctx{
		Notify: notifications.Notify,
		GetCsvUrl: func() string {
			return os.Getenv("DOMAINS_CSV_URL")
		},
	}

	ctx.Run()
}

func main() {
	lambda.Start(Handler)
}
