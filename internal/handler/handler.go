package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/fre5h/rnokpp"
	"log"
	"net/http"
	"strconv"

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

	if err := json.Unmarshal([]byte(request.Body), &update); err != nil {
		log.Print("error on unmarshal request body: ", err)

		return createLambdaResponse(http.StatusInternalServerError, "Error on unmarshal json")
	}

	if 0 == update.UpdateId {
		return createLambdaResponse(http.StatusBadRequest, "Update id of 0 indicates failure to parse incoming update")
	}

	var sendMessage *model.SendMessage
	var chatId int

	if nil != update.Message && "" != update.Message.Text {
		chatId = update.Message.Chat.Id

		if nil != update.Message.ReplyToMessage && update.Message.ReplyToMessage.Text == "Введіть РНОКПП по якому ви хочете отримати інформацію:" {
			sendMessage = processGetInfo(update.Message.Text)
		} else {
			sendMessage = prepareResult(update.Message.Text)
		}
	} else if nil != update.CallbackQuery {
		chatId = update.CallbackQuery.Message.Chat.Id
		sendMessage = prepareResult(*update.CallbackQuery.Data)
	}

	if nil != sendMessage && 0 != chatId {
		if responseBody, err := h.telegramClient.SendMessageToChat(chatId, *sendMessage); err != nil {
			log.Printf("error %s from telegram, response body is %s", err.Error(), responseBody)

			return createLambdaResponse(http.StatusInternalServerError, "Error on request to Telegram")
		}
	}

	return createLambdaResponse(http.StatusOK, "OK")
}

func createLambdaResponse(statusCode int, body string) (events.LambdaFunctionURLResponse, error) {
	return events.LambdaFunctionURLResponse{StatusCode: statusCode, Body: body}, nil
}

func prepareResult(text string) *model.SendMessage {
	switch text {
	case "/start":
		return &model.SendMessage{
			Text: "Просто напишіть РНОКПП, який ви хочете розшифрувати. Або виберіть команду з меню. Або введіть символ / і виберіть потрібну команду.",
		}
	case "/random":
		randomRnokpp, _ := rnokpp.GenerateRandomRnokpp()
		details, _ := rnokpp.GetDetails(randomRnokpp)

		return &model.SendMessage{
			Text: fmt.Sprintf(
				"Випадковий РНОКПП: \n<code>%s</code> - стать: <code>%s</code>, дата народження: <code>%s</code>",
				randomRnokpp,
				ukrGender(details.Gender),
				details.Birthday.Format("02.01.2006"),
			),
			ParseMode: "HTML",
		}
	case "/randomn":
		return &model.SendMessage{
			Text: "Виберіть кількість випадкових РНОКПП, яку ви хочете згенерувати:",
			ReplyMarkup: &model.ReplyMarkup{
				OneTimeKeyboard: true,
				IsPersistent:    true,
				ResizeKeyboard:  true,
				InlineKeyboard: [][]model.InlineKeyboardButton{
					{
						model.InlineKeyboardButton{Text: "2", CallbackData: "random2"},
						model.InlineKeyboardButton{Text: "3", CallbackData: "random3"},
						model.InlineKeyboardButton{Text: "5", CallbackData: "random5"},
						model.InlineKeyboardButton{Text: "10", CallbackData: "random10"},
						model.InlineKeyboardButton{Text: "20", CallbackData: "random20"},
					},
				},
			},
		}
	case "random2", "random3", "random5", "random10", "random20":
		n, _ := strconv.Atoi(text[6:])
		if n > 0 {
			return processRandomN(n)
		}
		return nil
	case "/getinfo":
		return &model.SendMessage{
			Text:                     "Введіть РНОКПП по якому ви хочете отримати інформацію:",
			AllowSendingWithoutReply: false,
			ReplyMarkup: &model.ReplyMarkup{
				ForceReply: true,
			},
		}
	default:
		return processGetInfo(text)
	}
}

func processRandomN(number int) *model.SendMessage {
	result := "Випадкові РНОКПП:"

	for i := 0; i < number; i++ {
		randomRnokpp, _ := rnokpp.GenerateRandomRnokpp()
		details, _ := rnokpp.GetDetails(randomRnokpp)

		result += fmt.Sprintf(
			"\n<code>%s</code> - стать: <code>%s</code>, дата народження: <code>%s</code>",
			randomRnokpp,
			ukrGender(details.Gender),
			details.Birthday.Format("02.01.2006"),
		)
	}

	return &model.SendMessage{
		Text:        result,
		ParseMode:   "HTML",
		ReplyMarkup: nil,
	}
}

func processGetInfo(text string) *model.SendMessage {
	details, err := rnokpp.GetDetails(text)

	if errors.Is(err, rnokpp.ErrMoreThan10Digits) ||
		errors.Is(err, rnokpp.ErrLessThan10Digits) ||
		errors.Is(err, rnokpp.ErrStringDoesNotConsistOfDigits) {
		return &model.SendMessage{
			Text: "Невірний формат РНОКПП. Повинно бути рівно 10 цифр.",
		}
	}

	if errors.Is(err, rnokpp.ErrInvalidControlDigit) {
		return &model.SendMessage{
			Text: "Введений текст не є валідним РНОКПП. Валідний РНОКПП повинен складатись з 10 цифр, містити в собі зашифровану інформацію про стать, дату народження та контрольну суму.",
		}
	}

	return &model.SendMessage{
		ParseMode: "HTML",
		Text: fmt.Sprintf(
			"РНОКПП: <code>%s</code>\nСтать: <code>%s</code>\nДата народження: <code>%s</code>",
			text,
			ukrGender(details.Gender),
			details.Birthday.Format("02.01.2006"),
		),
	}
}

func ukrGender(gender rnokpp.Gender) string {
	switch gender {
	case rnokpp.Male:
		return "чоловіча"
	case rnokpp.Female:
		return "жіноча"
	default:
		panic("unknown gender")
	}
}
