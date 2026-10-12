package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fre5h/rnokpp-telegram-bot/internal/model"
)

const (
	defaultBaseURL = "https://api.telegram.org/bot"
	requestTimeout = 10 * time.Second
)

// Client sends messages to Telegram chats.
type Client interface {
	// SendMessageToChat returns the raw response body. It is returned together with an error
	// when Telegram rejects the request, so the caller can log Telegram's explanation.
	SendMessageToChat(chatID int, message model.SendMessage) (string, error)
}

type HTTPClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewHTTPClient creates a client for the Telegram Bot API.
func NewHTTPClient(token string) (*HTTPClient, error) {
	return newHTTPClient(defaultBaseURL, token, &http.Client{Timeout: requestTimeout})
}

func newHTTPClient(baseURL, token string, httpClient *http.Client) (*HTTPClient, error) {
	if token == "" {
		return nil, errors.New("telegram bot token is empty")
	}

	return &HTTPClient{baseURL: baseURL, token: token, httpClient: httpClient}, nil
}

func (c *HTTPClient) SendMessageToChat(chatID int, message model.SendMessage) (string, error) {
	message.ChatID = strconv.Itoa(chatID)

	payload, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("marshal message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+c.token+"/sendMessage", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("build request: %w", withoutURL(err))
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("post message to chat: %w", withoutURL(err))
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("error on closing body: %s", err.Error())
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read telegram response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return string(body), fmt.Errorf("telegram responded with status %d", response.StatusCode)
	}

	return string(body), nil
}

// withoutURL drops the *url.Error wrapper, whose message contains the request URL and therefore
// the bot token. The underlying cause is kept.
func withoutURL(err error) error {
	if urlErr, ok := errors.AsType[*url.Error](err); ok {
		return urlErr.Err
	}

	return err
}
