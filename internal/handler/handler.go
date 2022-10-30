package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/fre5h/rnokpp-telegram-bot/internal/model"
)

type LambdaHandler struct {
	telegramClient TelegramClient
}

func NewLambdaHandler(telegramClient TelegramClient) *LambdaHandler {
	return &LambdaHandler{telegramClient: telegramClient}
}

func (h LambdaHandler) HandleLambdaRequest(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	var update model.Update

	err := json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		return createLambdaResponse(http.StatusInternalServerError, "Error on unmarshal json")
	}

	if 0 == update.UpdateId {
		return createLambdaResponse(http.StatusBadRequest, "Update id of 0 indicates failure to parse incoming update")
	}

	if responseBody, err := h.telegramClient.SendTextMessageToChat(update.Message.Chat.Id, prepareResult(update.Message.Text)); err != nil {
		log.Printf("error %s from telegram, response body is %s", err.Error(), responseBody)

		return createLambdaResponse(http.StatusInternalServerError, "Error on request to Telegram")
	}

	return createLambdaResponse(http.StatusOK, "OK")
}

func createLambdaResponse(statusCode int, body string) (events.LambdaFunctionURLResponse, error) {
	return events.LambdaFunctionURLResponse{StatusCode: statusCode, Body: body}, nil
}

func prepareResult(text string) (result string) {
	switch text {
	case "":
		result = "🤔 Вибачайте, але я вмію транслітерувати лише текстові повідомлення"
	case "/start":
		result = "Просто напишіть мені текст на українській мові 🇺🇦 і у відповідь отримаєте транслітерований 🇬🇧 текст"
	default:
		result = "123"
	}

	return result
}
