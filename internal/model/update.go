package model

type Chat struct {
	Id int `json:"id"`
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
	Id              string   `json:"id"`
	Data            *string  `json:"data,omitempty"`
	Message         *Message `json:"message,omitempty"`
	InlineMessageId *string  `json:"inline_message_id,omitempty"`
}

type Update struct {
	UpdateId      int            `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}
