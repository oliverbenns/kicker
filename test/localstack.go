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
	log.Print("waiting for localstack to start")
	deadline := time.Now().Add(20 * time.Second)

	for {
		if deadline.Before(time.Now()) {
			return errors.New("timed out waiting for localstack to start")
		}

		time.Sleep(1 * time.Second)

		r, err := http.Get("http://localhost:4566")
		if err != nil {
			continue
		}

		// cannot just check for 200...returns 404 but with a res body
		var status LocalStackStatus
		err = json.NewDecoder(r.Body).Decode(&status)
		if err != nil || status.Status != "running" {
			continue
		}

		break
	}

	return nil
}
