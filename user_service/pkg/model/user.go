package model

import "time"

var UserTable = "users"

type User struct {
	ID           uint64     `json:"id" sql:"id"`
	Email        string     `json:"email" db:"email"`
	Name         string     `json:"name" db:"name"`
	PasswordHash string     `json:"password_hash" db:"password_hash"`
	Description  string     `json:"description" db:"description"`
	Avatar       string     `json:"avatar" db:"avatar"`
	CreatedAt    *time.Time `json:"created_at" db:"crated_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
	Role         Role       `json:"role"`
}
