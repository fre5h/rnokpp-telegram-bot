BUILD_ENVS = GOARCH=arm64 GOOS=linux CGO_ENABLED=0

.PHONY: aws_lambda test

aws_lambda:
	${BUILD_ENVS} go build -trimpath -ldflags="-s -w" -tags lambda.norpc -o bootstrap cmd/main.go
	rm -f archive.zip
	zip archive.zip bootstrap

test:
	go test -race ./...
