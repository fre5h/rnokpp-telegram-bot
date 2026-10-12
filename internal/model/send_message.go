package model

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

// InlineKeyboardMarkup is the reply_markup for buttons attached to a message.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// ForceReplyMarkup asks the user to reply to the message, which the bot uses to detect input.
type ForceReplyMarkup struct {
	ForceReply bool `json:"force_reply"`
}

type SendMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
	// ReplyMarkup is one of *InlineKeyboardMarkup or *ForceReplyMarkup, or nil.
	ReplyMarkup any `json:"reply_markup,omitempty"`
}
