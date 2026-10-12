package handler

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/fre5h/rnokpp-telegram-bot/internal/bot"
	"github.com/fre5h/rnokpp-telegram-bot/internal/model"
	"github.com/fre5h/rnokpp-telegram-bot/internal/telegram"
)

// secretTokenHeader is sent by Telegram when the webhook was registered with a secret_token.
const secretTokenHeader = "X-Telegram-Bot-Api-Secret-Token"

// LambdaHandler adapts Telegram webhook requests to the bot.
type LambdaHandler struct {
	bot            *bot.Bot
	telegramClient telegram.Client
	// webhookSecret is checked against secretTokenHeader. Empty means requests are not authenticated.
	webhookSecret string
}

func NewLambdaHandler(bot *bot.Bot, telegramClient telegram.Client, webhookSecret string) *LambdaHandler {
	return &LambdaHandler{bot: bot, telegramClient: telegramClient, webhookSecret: webhookSecret}
}

func (h *LambdaHandler) HandleLambdaRequest(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	if !h.authorized(request.Headers) {
		return createLambdaResponse(http.StatusUnauthorized, "Unauthorized")
	}

	body, err := requestBody(request)
	if err != nil {
		log.Print("error on decoding request body: ", err)

		return createLambdaResponse(http.StatusBadRequest, "Invalid request body")
	}

	var update model.Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Print("error on unmarshal request body: ", err)

		return createLambdaResponse(http.StatusBadRequest, "Error on unmarshal json")
	}

	if 0 == update.UpdateID {
		return createLambdaResponse(http.StatusBadRequest, "Update id of 0 indicates failure to parse incoming update")
	}

	chatID, reply := h.bot.Handle(update)
	if nil == reply || 0 == chatID {
		return createLambdaResponse(http.StatusOK, "OK")
	}

	// Telegram redelivers the update on any non-2xx response, so a failed send is logged and
	// acknowledged instead of failing the request. Otherwise the same update would be processed again.
	if responseBody, err := h.telegramClient.SendMessageToChat(chatID, *reply); err != nil {
		log.Printf("error on sending message to chat %d: %v, telegram response: %s", chatID, err, responseBody)
	}

	return createLambdaResponse(http.StatusOK, "OK")
}

func (h *LambdaHandler) authorized(headers map[string]string) bool {
	if "" == h.webhookSecret {
		return true
	}

	for name, value := range headers {
		if strings.EqualFold(name, secretTokenHeader) {
			return 1 == subtle.ConstantTimeCompare([]byte(value), []byte(h.webhookSecret))
		}
	}

	return false
}

func requestBody(request events.LambdaFunctionURLRequest) ([]byte, error) {
	if request.IsBase64Encoded {
		return base64.StdEncoding.DecodeString(request.Body)
	}

	return []byte(request.Body), nil
}

func createLambdaResponse(statusCode int, body string) (events.LambdaFunctionURLResponse, error) {
	return events.LambdaFunctionURLResponse{StatusCode: statusCode, Body: body}, nil
}
