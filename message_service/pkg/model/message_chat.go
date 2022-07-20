package model

import "time"

const (
	MESSAGE_CHAT_TABLE = "messages_chats"
)

type MessageChat struct {
	ID           int        `json:"id" sql:"id"`
	ChatId       int        `json:"chat_id" db:"chat_id"`
	SenderUserId int        `json:"sender_user_id" db:"sender_user_id"`
	Body         string     `json:"body" db:"body"`
	CreatedAt    *time.Time `json:"created_at" db:"crated_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
}
type FilterMessage struct {
	UserId int `json:"user_id"`
	ChatId int `json:"chat_id"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
