package test

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type LocalStackStatus struct {
	Status string
}

func WaitForLocalStack() error {
	deadline := time.Now().Add(20 * time.Second)
	printed := false

	for {
		if deadline.Before(time.Now()) {
			return errors.New("timed out waiting for localstack to start")
		}

		r, err := http.Get("http://localhost:4566")
		if err != nil {
			if !printed {
				printed = true
				log.Print("waiting for localstack to start")

			}
			continue
		}

		// cannot just check for 200...returns 404 but with a res body
		var status LocalStackStatus
		err = json.NewDecoder(r.Body).Decode(&status)
		if err != nil || status.Status != "running" {
			if !printed {
				printed = true
				log.Print("waiting for localstack to start")
			}
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	return nil
}
