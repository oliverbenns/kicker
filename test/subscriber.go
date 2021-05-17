package test

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type notification struct {
	Message string
}
type subscriber struct {
	notifications chan notification
	errors        chan error
	sns           *testSns
}

func NewSubscriber(testSns *testSns) (*subscriber, error) {
	return &subscriber{
		errors:        make(chan error),
		notifications: make(chan notification),
		sns:           testSns,
	}, nil
}

func (s *subscriber) Listen() error {
	_, err := s.sns.Client.Subscribe(&sns.SubscribeInput{
		Endpoint: aws.String("http://host.docker.internal:5000"),
		Protocol: aws.String("http"),
		TopicArn: s.sns.TopicArn,
	})
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var n notification
		err := json.NewDecoder(r.Body).Decode(&n)
		if err != nil {
			s.errors <- err
		}

		s.notifications <- n
	})

	return http.ListenAndServe(":5000", nil)
}
func (s *subscriber) GetNotification() (*notification, error) {
	select {
	case n := <-s.notifications:
		return &n, nil
	case err := <-s.errors:
		return nil, err
	case <-time.After(5 * time.Second):
		return nil, errors.New("timeout waiting for notification")
	}
}
