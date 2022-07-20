package model

const USER_UNREAD_MESSAGE = "users_unread_messages"

type UserUnreadMessage struct {
	ID        int `json:"id" sql:"id"`
	ChatId    int `json:"chat_id" db:"chat_id"`
	MessageId int `json:"message_id" db:"message_id"`
	UserId    int `json:"user_id" db:"user_id"`
}
