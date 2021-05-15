.PHONY: build clean deploy test

build:
	go build -o bin/domainfree cmd/domainfree/main.go
	go build -o bin/ping cmd/ping/main.go

clean:
	rm -rf ./bin
	rm -rf ./zip

test:
	go test cmd/domainfree/*.go -v
	go test cmd/ping/*.go -v

deploy: clean build
	mkdir zip
	zip zip/ping.zip bin/ping
	zip zip/domainfree.zip bin/domainfree
	aws s3 sync zip s3://${AWS_S3_BUCKET_NAME} --delete
