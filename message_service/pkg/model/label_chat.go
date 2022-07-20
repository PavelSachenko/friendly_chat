package model

const LABELS_CHATS_TABLE = "labels_chats"

type LabelChats struct {
	ID     int    `json:"id" sql:"id"`
	Title  string `json:"title" db:"title"`
	UserId int    `json:"user_id" db:"user_id"`
	ChatId int    `json:"chat_id" db:"chat_id"`
}

type UserTitleChats struct {
	UserId int    `json:"user_id"`
	Title  string `json:"title"`
}
