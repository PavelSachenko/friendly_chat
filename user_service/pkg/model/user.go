package model

import "time"

var UserTable = "users"

type User struct {
	ID           uint64     `json:"id" sql:"id"`
	Username     string     `json:"username" db:"username"`
	PasswordHash string     `json:"password_hash" db:"password_hash"`
	Description  string     `json:"description" db:"description"`
	Avatar       string     `json:"avatar" db:"avatar"`
	CreatedAt    *time.Time `json:"created_at" db:"crated_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
	Role         Role       `json:"role"`
}
