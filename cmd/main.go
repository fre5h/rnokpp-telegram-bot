package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fre5h/rnokpp-telegram-bot/internal/handler"
)

func main() {
	httpClient := &http.Client{}
	telegramClient := handler.NewTelegramHttpClient(httpClient)
	lambdaHandler := handler.NewLambdaHandler(*telegramClient)
	lambda.Start(lambdaHandler.HandleLambdaRequest)
}
