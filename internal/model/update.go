package model

type Chat struct {
	ID int `json:"id"`
}

type ReplyToMessage struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Message struct {
	Text           string          `json:"text"`
	Chat           Chat            `json:"chat"`
	ReplyToMessage *ReplyToMessage `json:"reply_to_message,omitempty"`
}

type CallbackQuery struct {
	ID              string   `json:"id"`
	Data            *string  `json:"data,omitempty"`
	Message         *Message `json:"message,omitempty"`
	InlineMessageID *string  `json:"inline_message_id,omitempty"`
}

type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}
