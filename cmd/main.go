package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fre5h/rnokpp-telegram-bot/internal/handler"
)

func main() {
	telegramClient := handler.NewTelegramHttpClient()
	lambdaHandler := handler.NewLambdaHandler(*telegramClient)
	lambda.Start(lambdaHandler.HandleLambdaRequest)
}
