BUILD_ENVS = GOARCH=arm64 GOOS=linux

aws_lambda:
	${BUILD_ENVS} go build -tags lambda.norpc -o bootstrap cmd/main.go && zip archive.zip bootstrap
