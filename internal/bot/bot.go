package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/fre5h/rnokpp"

	"github.com/fre5h/rnokpp-telegram-bot/internal/model"
)

// PromptRnokpp is the prompt shown with ForceReply. A reply to it is treated as an RNOKPP lookup.
const PromptRnokpp = "Введіть РНОКПП по якому ви хочете отримати інформацію:"

const (
	commandStart   = "/start"
	commandRandom  = "/random"
	commandRandomN = "/randomn"
	commandGetInfo = "/getinfo"
)

const (
	textStart             = "Просто напишіть РНОКПП, який ви хочете розшифрувати. Або виберіть команду з меню. Або введіть символ / і виберіть потрібну команду."
	textRandomCountPrompt = "Виберіть кількість випадкових РНОКПП, яку ви хочете згенерувати:"
	textInvalidFormat     = "Невірний формат РНОКПП. Повинно бути рівно 10 цифр."
	textInvalidRnokpp     = "Введений текст не є валідним РНОКПП. Валідний РНОКПП повинен складатись з 10 цифр, містити в собі зашифровану інформацію про стать, дату народження та контрольну суму."
	textGenerateFailed    = "Не вдалося згенерувати РНОКПП. Спробуйте ще раз."
	headerRandomSingle    = "Випадковий РНОКПП:"
	headerRandomMany      = "Випадкові РНОКПП:"
)

// randomCountButtons are the inline keyboard buttons of /randomn, in display order.
var randomCountButtons = []struct {
	callbackData string
	count        int
}{
	{"random2", 2},
	{"random3", 3},
	{"random5", 5},
	{"random10", 10},
	{"random20", 20},
}

// Generator returns count random valid RNOKPPs. rnokpp.GenerateRandomRnokppN fits this signature.
type Generator func(count int) ([]string, error)

type Bot struct {
	generate Generator
}

func New(generate Generator) *Bot {
	return &Bot{generate: generate}
}

// Handle returns the chat to reply to and the reply. A nil reply means nothing should be sent.
func (b *Bot) Handle(update model.Update) (int, *model.SendMessage) {
	if query := update.CallbackQuery; query != nil {
		// Inline-mode queries have no Message, and any query may lack Data.
		if query.Message == nil || query.Data == nil {
			return 0, nil
		}

		return query.Message.Chat.ID, b.handleCallback(*query.Data)
	}

	if message := update.Message; message != nil && message.Text != "" {
		if message.ReplyToMessage != nil && message.ReplyToMessage.Text == PromptRnokpp {
			return message.Chat.ID, processGetInfo(message.Text)
		}

		return message.Chat.ID, b.handleCommand(message.Text)
	}

	return 0, nil
}

func (b *Bot) handleCallback(data string) *model.SendMessage {
	for _, button := range randomCountButtons {
		if button.callbackData == data {
			return b.randomReply(headerRandomMany, button.count)
		}
	}

	return nil
}

func (b *Bot) handleCommand(text string) *model.SendMessage {
	switch text {
	case commandStart:
		return &model.SendMessage{Text: textStart}
	case commandRandom:
		return b.randomReply(headerRandomSingle, 1)
	case commandRandomN:
		return &model.SendMessage{Text: textRandomCountPrompt, ReplyMarkup: randomCountKeyboard()}
	case commandGetInfo:
		return &model.SendMessage{Text: PromptRnokpp, ReplyMarkup: &model.ForceReplyMarkup{ForceReply: true}}
	default:
		return processGetInfo(text)
	}
}

func randomCountKeyboard() *model.InlineKeyboardMarkup {
	row := make([]model.InlineKeyboardButton, 0, len(randomCountButtons))
	for _, button := range randomCountButtons {
		row = append(row, model.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", button.count),
			CallbackData: button.callbackData,
		})
	}

	return &model.InlineKeyboardMarkup{InlineKeyboard: [][]model.InlineKeyboardButton{row}}
}

func (b *Bot) randomReply(header string, count int) *model.SendMessage {
	lines, err := b.randomLines(count)
	if err != nil {
		log.Printf("error on generating random rnokpp: %v", err)

		return &model.SendMessage{Text: textGenerateFailed}
	}

	return &model.SendMessage{
		Text:      header + "\n" + strings.Join(lines, "\n"),
		ParseMode: "HTML",
	}
}

func (b *Bot) randomLines(count int) ([]string, error) {
	values, err := b.generate(count)
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}

	lines := make([]string, 0, len(values))
	for _, value := range values {
		details, err := rnokpp.GetDetails(value)
		if err != nil {
			return nil, fmt.Errorf("details of generated rnokpp: %w", err)
		}

		lines = append(lines, formatRandomLine(value, details))
	}

	return lines, nil
}

func formatRandomLine(value string, details *rnokpp.Details) string {
	return fmt.Sprintf(
		"<code>%s</code> - стать: <code>%s</code>, дата народження: <code>%s</code>",
		value,
		ukrainianGender(details.Gender),
		details.Birthday.Format("02.01.2006"),
	)
}

func processGetInfo(text string) *model.SendMessage {
	details, err := rnokpp.GetDetails(text)

	switch {
	case errors.Is(err, rnokpp.ErrMoreThan10Digits),
		errors.Is(err, rnokpp.ErrLessThan10Digits),
		errors.Is(err, rnokpp.ErrStringDoesNotConsistOfDigits):
		return &model.SendMessage{Text: textInvalidFormat}
	case errors.Is(err, rnokpp.ErrInvalidControlDigit):
		return &model.SendMessage{Text: textInvalidRnokpp}
	case err != nil:
		log.Printf("unexpected error on getinfo: %v", err)

		return &model.SendMessage{Text: textInvalidRnokpp}
	}

	return &model.SendMessage{
		ParseMode: "HTML",
		Text: fmt.Sprintf(
			"РНОКПП: <code>%s</code>\nСтать: <code>%s</code>\nДата народження: <code>%s</code>",
			text,
			ukrainianGender(details.Gender),
			details.Birthday.Format("02.01.2006"),
		),
	}
}

func ukrainianGender(gender rnokpp.Gender) string {
	switch gender {
	case rnokpp.Male:
		return "чоловіча"
	case rnokpp.Female:
		return "жіноча"
	default:
		return "невідома"
	}
}
