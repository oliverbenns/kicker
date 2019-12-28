include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build clean deploy test

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/check-domain-availability check-domain-availability/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/check-website-up check-website-up/main.go

clean:
	rm -rf ./bin

test:
	go test check-domain-availability/*.go -v
	go test check-website-up/*.go -v

deploy: clean build
	sls deploy --verbose
