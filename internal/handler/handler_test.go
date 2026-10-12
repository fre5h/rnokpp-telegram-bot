package handler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/fre5h/rnokpp-telegram-bot/internal/bot"
	"github.com/fre5h/rnokpp-telegram-bot/internal/model"
)

type fakeClient struct {
	calls   int
	chatID  int
	text    string
	sendErr error
}

func (f *fakeClient) SendMessageToChat(chatID int, message model.SendMessage) (string, error) {
	f.calls++
	f.chatID = chatID
	f.text = message.Text

	return `{"ok":true}`, f.sendErr
}

func fakeGenerator(count int) ([]string, error) {
	return make([]string, count), nil
}

const startUpdate = `{"update_id":1,"message":{"text":"/start","chat":{"id":42}}}`

func newHandler(client *fakeClient, secret string) *LambdaHandler {
	return NewLambdaHandler(bot.New(fakeGenerator), client, secret)
}

func request(body string) events.LambdaFunctionURLRequest {
	return events.LambdaFunctionURLRequest{Body: body}
}

func TestHandleSendsReply(t *testing.T) {
	client := &fakeClient{}

	response, err := newHandler(client, "").HandleLambdaRequest(request(startUpdate))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", response.StatusCode)
	}
	if client.calls != 1 || client.chatID != 42 {
		t.Errorf("sent %d messages to chat %d, want 1 to chat 42", client.calls, client.chatID)
	}
}

func TestHandleBadRequests(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "malformed json", body: `{"update_id":`},
		{name: "update id zero", body: `{"update_id":0,"message":{"text":"/start","chat":{"id":42}}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &fakeClient{}

			response, _ := newHandler(client, "").HandleLambdaRequest(request(tt.body))

			if response.StatusCode != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", response.StatusCode)
			}
			if client.calls != 0 {
				t.Errorf("sent %d messages, want 0", client.calls)
			}
		})
	}
}

func TestHandleSecretToken(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		wantStatus int
	}{
		{name: "missing header", headers: nil, wantStatus: http.StatusUnauthorized},
		{name: "wrong secret", headers: map[string]string{secretTokenHeader: "wrong"}, wantStatus: http.StatusUnauthorized},
		{name: "correct secret", headers: map[string]string{secretTokenHeader: "s3cret"}, wantStatus: http.StatusOK},
		{name: "header name is case insensitive", headers: map[string]string{"x-telegram-bot-api-secret-token": "s3cret"}, wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &fakeClient{}
			req := request(startUpdate)
			req.Headers = tt.headers

			response, _ := newHandler(client, "s3cret").HandleLambdaRequest(req)

			if response.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", response.StatusCode, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusUnauthorized && client.calls != 0 {
				t.Errorf("sent %d messages for unauthorized request, want 0", client.calls)
			}
		})
	}
}

func TestHandleWithoutSecretAcceptsAnyRequest(t *testing.T) {
	client := &fakeClient{}

	response, _ := newHandler(client, "").HandleLambdaRequest(request(startUpdate))

	if response.StatusCode != http.StatusOK || client.calls != 1 {
		t.Errorf("status = %d, sent = %d, want 200 and 1 message", response.StatusCode, client.calls)
	}
}

func TestHandleBase64Body(t *testing.T) {
	client := &fakeClient{}
	req := events.LambdaFunctionURLRequest{
		Body:            base64.StdEncoding.EncodeToString([]byte(startUpdate)),
		IsBase64Encoded: true,
	}

	response, _ := newHandler(client, "").HandleLambdaRequest(req)

	if response.StatusCode != http.StatusOK || client.calls != 1 {
		t.Errorf("status = %d, sent = %d, want 200 and 1 message", response.StatusCode, client.calls)
	}
}

func TestHandleAcknowledgesFailedSend(t *testing.T) {
	client := &fakeClient{sendErr: errors.New("telegram rejected message")}

	response, err := newHandler(client, "").HandleLambdaRequest(request(startUpdate))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// A non-2xx here would make Telegram redeliver the same update.
	if response.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", response.StatusCode)
	}
}

func TestHandleUpdateWithoutReplySendsNothing(t *testing.T) {
	client := &fakeClient{}

	response, _ := newHandler(client, "").HandleLambdaRequest(request(`{"update_id":1}`))

	if response.StatusCode != http.StatusOK || client.calls != 0 {
		t.Errorf("status = %d, sent = %d, want 200 and 0 messages", response.StatusCode, client.calls)
	}
}
