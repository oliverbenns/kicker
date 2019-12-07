.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/check-domain-availability check-domain-availability/main.go

clean:
	rm -rf ./bin

test:
	go test check-domain-availability/*.go -v

deploy: clean build
	sls deploy --verbose
