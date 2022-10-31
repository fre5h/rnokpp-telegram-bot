package model

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type ReplyMarkup struct {
	InlineKeyboard        [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
	OneTimeKeyboard       bool                     `json:"one_time_keyboard,omitempty"`
	IsPersistent          bool                     `json:"is_persistent,omitempty"`
	ResizeKeyboard        bool                     `json:"resize_keyboard,omitempty"`
	ForceReply            bool                     `json:"force_reply,omitempty"`
	InputFieldPlaceholder string                   `json:"input_field_placeholder,omitempty"`
}

type SendMessage struct {
	ChatId                   string       `json:"chat_id"`
	Text                     string       `json:"text"`
	ParseMode                string       `json:"parse_mode,omitempty"`
	AllowSendingWithoutReply bool         `json:"allow_sending_without_reply"`
	ReplyMarkup              *ReplyMarkup `json:"reply_markup,omitempty"`
}
