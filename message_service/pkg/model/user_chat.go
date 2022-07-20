package model

const USERS_CHATS_TABLE = "users_chats"

type UserChat struct {
	ID      int  `json:"id" sql:"id"`
	UserId  int  `json:"user_id" db:"user_id"`
	ChatId  int  `json:"chat_id" db:"chat_id"`
	IsOwner bool `json:"is_owner" db:"is_owner"`
}

type UsersOwners struct {
	UserId  int  `json:"user_id"`
	IsOwner bool `json:"is_owner"`
}
