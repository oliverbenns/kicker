.PHONY: build clean deploy test

build:
	go build -o bin/domainfree cmd/domainfree/main.go
	go build -o bin/websiteup cmd/websiteup/main.go

clean:
	rm -rf ./bin

test:
	go test cmd/domainfree/*.go -v
	go test cmd/websiteup/*.go -v

deploy: clean build
	sls deploy --verbose
