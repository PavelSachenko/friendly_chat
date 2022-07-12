package model

var UserRoleTable = "users_roles"

type UserRole struct {
	ID     uint64 `json:"id" sql:"id" db:"id"`
	UserId uint64 `json:"user_id" db:"user_id"`
	RoleId uint64 `json:"role_id" db:"role_id"`
}
