package model

import "time"

const (
	CHATS_TABLE  = "chats"
	PRIVATE_CHAT = 1
	GROUP_CHAT   = 2
)

type Chat struct {
	ID          int        `json:"id" sql:"id"`
	Description string     `json:"description" db:"description"`
	CreatedAt   *time.Time `json:"created_at" db:"crated_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
}

type AllChat struct {
	ID                 int          `json:"id" sql:"id"`
	Description        string       `json:"description" db:"description"`
	Title              string       `json:"title" db:"title"`
	CountUnreadMessage int          `json:"count_unread_message" db:"count_unread_message"`
	Type               int          `json:"type" db:"type"`
	StartPosition      int          `json:"start_position" db:"start_position"`
	LastMessage        *LastMessage `json:"last_message"`
	CreatedAt          *time.Time   `json:"created_at" db:"created_at"`
	UpdateAt           *time.Time   `json:"updated_at" db:"updated_at"`
}

type LastMessage struct {
	Body      string `json:"body"`
	IsSent    bool   `json:"is_sent"`
	CreatedAt string `json:"created_at"`
}

type FilterChat struct {
	UserId int `json:"user_id"`
	Limit  int `json:"limit" default:"20"`
	Offset int `json:"offset" default:"0"`
}
