.PHONY: build clean deploy test

build:
	go build -o bin/domainfree cmd/domainfree/main.go
	go build -o bin/ping cmd/ping/main.go

clean:
	rm -rf ./bin
	rm -rf ./zip

start-localstack:
	docker run -d --rm -it -p 4566:4566 -p 4571:4571 localstack/localstack

test:
	go test cmd/domainfree/*.go -v
	go test cmd/ping/*.go -v

deploy: clean build
	mkdir zip
	zip zip/ping.zip bin/ping
	zip zip/domainfree.zip bin/domainfree
	aws s3 sync zip s3://${AWS_S3_BUCKET_NAME} --delete
