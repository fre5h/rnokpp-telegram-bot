package handler

import (
	"errors"
	"testing"

	"github.com/fre5h/rnokpp-telegram-bot/internal/mocks"
)

func TestSendTextMessageToChatSuccessfully(t *testing.T) {
	client := NewTelegramHttpClient(mocks.NewMockHttpClient(200, "OK", nil))
	response, err := client.SendTextMessageToChat(1, "test")

	if err != nil {
		t.Errorf("Expected no error, but \"%s\" give", err)
	}

	if response != "OK" {
		t.Errorf("Expected response to be \"OK\", got %s", response)
	}
}

func TestSendTextMessageToChatWithErrorOnPost(t *testing.T) {
	client := NewTelegramHttpClient(mocks.NewMockHttpClient(200, "OK", errors.New("test")))
	response, err := client.SendTextMessageToChat(1, "test")

	if err == nil {
		t.Error("Expected to get an error")
	}

	if response != "" {
		t.Errorf("Expected response to be empty string, got %s", response)
	}
}

func TestSendTextMessageToChatWithNon200Code(t *testing.T) {
	client := NewTelegramHttpClient(mocks.NewMockHttpClient(400, "Bad Request", nil))
	response, err := client.SendTextMessageToChat(1, "test")

	if err.Error() != "status code of response is: 400" {
		t.Error("Error message is different from expected")
	}

	if response != "" {
		t.Errorf("Expected response to be empty string, got %s", response)
	}
}

func TestSendTextMessageToChatWithErrorOnReadCloser(t *testing.T) {
	client := NewTelegramHttpClient(mocks.NewMockHttpClientFailedCloser(200, "OK", nil))
	response, err := client.SendTextMessageToChat(1, "test")

	if err == nil {
		t.Error("Expected to get an error")
	}

	if response != "" {
		t.Errorf("Expected response to be empty string, got %s", response)
	}
}
