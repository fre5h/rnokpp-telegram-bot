BUILD_ENVS = GOARCH=amd64 GOOS=linux

aws_lambda:
	${BUILD_ENVS} go build -o main cmd/main.go && zip archive.zip main