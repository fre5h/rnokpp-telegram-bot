package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fre5h/rnokpp-telegram-bot/internal/model"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type TelegramClient interface {
	SendMessageToChat(int, model.SendMessage) (string, error)
}

type TelegramHttpClient struct {
	baseUrl string
	token   string
}

func NewTelegramHttpClient() *TelegramHttpClient {
	return &TelegramHttpClient{
		baseUrl: "https://api.telegram.org/bot",
		token:   os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

func (c TelegramHttpClient) SendMessageToChat(charId int, sendMessage model.SendMessage) (string, error) {
	var botApiUrl = c.baseUrl + c.token + "/sendMessage"

	sendMessage.ChatId = strconv.Itoa(charId)

	payload, _ := json.Marshal(sendMessage)
	req, _ := http.NewRequest("POST", botApiUrl, bytes.NewBuffer(payload))
	req.Header.Add("content-type", "application/json")

	response, err := http.DefaultClient.Do(req)

	if nil != err {
		return "", fmt.Errorf("error when posting text to the chat: %s", err.Error())
	}

	if http.StatusOK != response.StatusCode {
		return "", fmt.Errorf("status code of response is: %d", response.StatusCode)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if nil != err {
			log.Printf("error on closing body: %s", err.Error())
		}
	}(response.Body)

	var bodyBytes, errRead = io.ReadAll(response.Body)

	if nil != errRead {
		return "", fmt.Errorf("error on parsing telegram response %s", errRead.Error())
	}

	return string(bodyBytes), nil
}
