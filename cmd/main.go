package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fre5h/rnokpp"

	"github.com/fre5h/rnokpp-telegram-bot/internal/bot"
	"github.com/fre5h/rnokpp-telegram-bot/internal/handler"
	"github.com/fre5h/rnokpp-telegram-bot/internal/telegram"
)

func main() {
	telegramClient, err := telegram.NewHTTPClient(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("error on creating telegram client: %v", err)
	}

	webhookSecret := os.Getenv("TELEGRAM_WEBHOOK_SECRET")
	if "" == webhookSecret {
		log.Print("TELEGRAM_WEBHOOK_SECRET is not set, webhook requests are not authenticated")
	}

	lambdaHandler := handler.NewLambdaHandler(bot.New(rnokpp.GenerateRandomRnokppN), telegramClient, webhookSecret)
	lambda.Start(lambdaHandler.HandleLambdaRequest)
}
