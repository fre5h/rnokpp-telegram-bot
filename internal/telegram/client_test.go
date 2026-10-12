package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fre5h/rnokpp-telegram-bot/internal/model"
)

const testToken = "123:SECRET-TOKEN"

func newTestClient(t *testing.T, baseURL string) *HTTPClient {
	t.Helper()

	client, err := newHTTPClient(baseURL, testToken, &http.Client{Timeout: time.Second})
	if err != nil {
		t.Fatalf("newHTTPClient: %v", err)
	}

	return client
}

func TestNewHTTPClientRejectsEmptyToken(t *testing.T) {
	if _, err := NewHTTPClient(""); err == nil {
		t.Error("expected an error for empty token")
	}
}

func TestSendMessageToChatSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/bot"+testToken+"/sendMessage" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("content-type = %s", got)
		}

		body, _ := io.ReadAll(r.Body)
		var sent model.SendMessage
		if err := json.Unmarshal(body, &sent); err != nil {
			t.Errorf("request body is not json: %v", err)
		}
		if sent.ChatID != "42" || sent.Text != "hello" {
			t.Errorf("sent = %+v", sent)
		}

		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL+"/bot")

	body, err := client.SendMessageToChat(42, model.SendMessage{Text: "hello"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body != `{"ok":true}` {
		t.Errorf("body = %s", body)
	}
}

func TestSendMessageToChatNon200ReturnsTelegramBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"ok":false,"description":"bot was blocked by the user"}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL+"/bot")

	body, err := client.SendMessageToChat(42, model.SendMessage{Text: "hello"})

	if err == nil || !strings.Contains(err.Error(), "403") {
		t.Errorf("error = %v, want status 403", err)
	}
	if !strings.Contains(body, "bot was blocked by the user") {
		t.Errorf("body = %q, want Telegram's description", body)
	}
}

func TestSendMessageToChatErrorDoesNotLeakToken(t *testing.T) {
	// Port 1 refuses connections, so the request fails inside the HTTP client.
	client := newTestClient(t, "http://127.0.0.1:1/bot")

	_, err := client.SendMessageToChat(42, model.SendMessage{Text: "hello"})

	if err == nil {
		t.Fatal("expected a connection error")
	}
	if strings.Contains(err.Error(), testToken) {
		t.Errorf("error leaks bot token: %v", err)
	}
}
