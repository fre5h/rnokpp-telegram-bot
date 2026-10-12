package bot

import (
	"errors"
	"strings"
	"testing"

	"github.com/fre5h/rnokpp"

	"github.com/fre5h/rnokpp-telegram-bot/internal/model"
)

// validRnokpp is generated once so the tests do not depend on a hardcoded checksum.
var validRnokpp = mustRandomRnokpp()

func mustRandomRnokpp() string {
	value, err := rnokpp.GenerateRandomRnokpp()
	if err != nil {
		panic(err)
	}

	return value
}

// invalidChecksumRnokpp changes the control digit of validRnokpp, which keeps the format valid.
func invalidChecksumRnokpp() string {
	last := (int(validRnokpp[9]-'0') + 1) % 10

	return validRnokpp[:9] + string(rune('0'+last))
}

func fakeGenerator(count int) ([]string, error) {
	values := make([]string, count)
	for i := range values {
		values[i] = validRnokpp
	}

	return values, nil
}

func textMessage(chatID int, text string) model.Update {
	return model.Update{
		UpdateID: 1,
		Message:  &model.Message{Text: text, Chat: model.Chat{ID: chatID}},
	}
}

func replyMessage(chatID int, text string) model.Update {
	return model.Update{
		UpdateID: 1,
		Message: &model.Message{
			Text:           text,
			Chat:           model.Chat{ID: chatID},
			ReplyToMessage: &model.ReplyToMessage{Text: PromptRnokpp},
		},
	}
}

func callbackQuery(chatID int, data string) model.Update {
	return model.Update{
		UpdateID: 1,
		CallbackQuery: &model.CallbackQuery{
			Data:    &data,
			Message: &model.Message{Chat: model.Chat{ID: chatID}},
		},
	}
}

func TestHandle(t *testing.T) {
	tests := []struct {
		name     string
		update   model.Update
		generate Generator
		wantChat int
		// wantNil is true when nothing should be sent.
		wantNil bool
		// wantText is a substring of the reply text.
		wantText string
	}{
		{
			name:     "start command",
			update:   textMessage(42, "/start"),
			generate: fakeGenerator,
			wantChat: 42,
			wantText: "Просто напишіть",
		},
		{
			name:     "random command generates one rnokpp",
			update:   textMessage(42, "/random"),
			generate: fakeGenerator,
			wantChat: 42,
			wantText: "Випадковий РНОКПП:",
		},
		{
			name:     "random count command shows keyboard",
			update:   textMessage(42, "/randomn"),
			generate: fakeGenerator,
			wantChat: 42,
			wantText: textRandomCountPrompt,
		},
		{
			name:     "getinfo command asks for rnokpp",
			update:   textMessage(42, "/getinfo"),
			generate: fakeGenerator,
			wantChat: 42,
			wantText: PromptRnokpp,
		},
		{
			name:     "callback generates requested count",
			update:   callbackQuery(42, "random3"),
			generate: fakeGenerator,
			wantChat: 42,
			wantText: headerRandomMany,
		},
		{
			name:     "unknown callback data sends nothing",
			update:   callbackQuery(42, "random999"),
			generate: fakeGenerator,
			wantChat: 42,
			wantNil:  true,
		},
		{
			name: "callback without message sends nothing",
			update: model.Update{
				UpdateID:      1,
				CallbackQuery: &model.CallbackQuery{Data: strPtr("random2")},
			},
			generate: fakeGenerator,
			wantNil:  true,
		},
		{
			name: "callback without data sends nothing",
			update: model.Update{
				UpdateID: 1,
				CallbackQuery: &model.CallbackQuery{
					Message: &model.Message{Chat: model.Chat{ID: 42}},
				},
			},
			generate: fakeGenerator,
			wantNil:  true,
		},
		{
			name:     "random-looking text is not a callback",
			update:   textMessage(42, "random5"),
			generate: fakeGenerator,
			wantChat: 42,
			wantText: textInvalidFormat,
		},
		{
			name:     "reply to prompt with valid rnokpp",
			update:   replyMessage(42, validRnokpp),
			generate: fakeGenerator,
			wantChat: 42,
			wantText: "Стать",
		},
		{
			name:     "reply to prompt with short rnokpp",
			update:   replyMessage(42, "123456789"),
			generate: fakeGenerator,
			wantChat: 42,
			wantText: textInvalidFormat,
		},
		{
			name:     "reply to prompt with wrong checksum",
			update:   replyMessage(42, invalidChecksumRnokpp()),
			generate: fakeGenerator,
			wantChat: 42,
			wantText: textInvalidRnokpp,
		},
		{
			name:     "generator error sends failure text",
			update:   textMessage(42, "/random"),
			generate: func(int) ([]string, error) { return nil, errors.New("boom") },
			wantChat: 42,
			wantText: textGenerateFailed,
		},
		{
			name:     "update without message sends nothing",
			update:   model.Update{UpdateID: 1},
			generate: fakeGenerator,
			wantNil:  true,
		},
		{
			name:     "empty text message sends nothing",
			update:   textMessage(42, ""),
			generate: fakeGenerator,
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatID, reply := New(tt.generate).Handle(tt.update)

			if tt.wantNil {
				if reply != nil {
					t.Fatalf("expected no reply, got %q", reply.Text)
				}

				return
			}

			if chatID != tt.wantChat {
				t.Errorf("chat id = %d, want %d", chatID, tt.wantChat)
			}

			if reply == nil {
				t.Fatal("expected a reply, got nil")
			}

			if !strings.Contains(reply.Text, tt.wantText) {
				t.Errorf("reply text %q does not contain %q", reply.Text, tt.wantText)
			}
		})
	}
}

func TestHandleRandomCountReturnsRequestedLines(t *testing.T) {
	_, reply := New(fakeGenerator).Handle(callbackQuery(42, "random5"))
	if reply == nil {
		t.Fatal("expected a reply, got nil")
	}

	// One newline after the header and one between each of the five lines.
	if got := strings.Count(reply.Text, "\n"); got != 5 {
		t.Errorf("got %d newlines in %q, want 5", got, reply.Text)
	}
}

func TestHandleRandomCountKeyboardHasButtonForEachCount(t *testing.T) {
	_, reply := New(fakeGenerator).Handle(textMessage(42, "/randomn"))

	markup, ok := reply.ReplyMarkup.(*model.InlineKeyboardMarkup)
	if !ok {
		t.Fatalf("reply markup is %T, want *model.InlineKeyboardMarkup", reply.ReplyMarkup)
	}

	if got := len(markup.InlineKeyboard[0]); got != len(randomCountButtons) {
		t.Errorf("got %d buttons, want %d", got, len(randomCountButtons))
	}

	for _, button := range markup.InlineKeyboard[0] {
		if _, _, found := strings.Cut(button.CallbackData, "random"); !found {
			t.Errorf("button %q has unexpected callback data %q", button.Text, button.CallbackData)
		}
	}
}

func TestProcessGetInfoRejectsChecksumWithoutPanicking(t *testing.T) {
	reply := processGetInfo(invalidChecksumRnokpp())

	if reply.Text != textInvalidRnokpp {
		t.Errorf("reply text = %q, want %q", reply.Text, textInvalidRnokpp)
	}
}

func TestUkrainianGender(t *testing.T) {
	if got := ukrainianGender(rnokpp.Male); got != "чоловіча" {
		t.Errorf("male = %q", got)
	}

	if got := ukrainianGender(rnokpp.Female); got != "жіноча" {
		t.Errorf("female = %q", got)
	}

	if got := ukrainianGender(rnokpp.Gender("unknown")); got != "невідома" {
		t.Errorf("unknown = %q, want невідома", got)
	}
}

func strPtr(s string) *string {
	return &s
}
