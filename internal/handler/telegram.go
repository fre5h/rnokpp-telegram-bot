package handler

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
)

type TelegramClient interface {
	SendTextMessageToChat(int, string) (string, error)
}

type TelegramHttpClient struct {
	baseUrl    string
	token      string
	httpClient HttpClient
}

func NewTelegramHttpClient(httpClient HttpClient) *TelegramHttpClient {
	return &TelegramHttpClient{
		baseUrl:    "https://api.telegram.org/bot",
		token:      os.Getenv("TELEGRAM_BOT_TOKEN"),
		httpClient: httpClient,
	}
}

func (c TelegramHttpClient) SendTextMessageToChat(chatId int, text string) (string, error) {
	var botApiUrl = c.baseUrl + c.token + "/sendMessage"

	response, err := c.httpClient.PostForm(
		botApiUrl,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		},
	)

	if nil != err {
		return "", fmt.Errorf("error when posting text to the chat: %s", err.Error())
	}

	if 200 != response.StatusCode {
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
		return "", fmt.Errorf("error on parsing telegram answer %s", errRead.Error())
	}

	return string(bodyBytes), nil
}
