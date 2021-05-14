.PHONY: build clean deploy test

build:
	go build -o bin/domainfree cmd/domainfree/main.go
	go build -o bin/ping cmd/ping/main.go

clean:
	rm -rf ./bin

test:
	go test cmd/domainfree/*.go -v
	go test cmd/ping/*.go -v

deploy: clean build
	sls deploy --verbose
