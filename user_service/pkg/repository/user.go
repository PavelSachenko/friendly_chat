package repository

import (
	"database/sql"
	"github.com/pavel/user_service/pkg/db"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/model"
)

type User interface {
	One(id uint64) (error, *model.User)
}

type UserPostgres struct {
	db     *db.DB
	logger *logger.Logger
}

func InitUserPostgres(db *db.DB, logger *logger.Logger) *UserPostgres {
	return &UserPostgres{
		db:     db,
		logger: logger,
	}
}

func (u UserPostgres) One(id uint64) (error, *model.User) {
	var user model.User
	rows, err := u.db.Queryx("SELECT id,email,name,description,avatar,created_at,updated_at FROM "+model.UserTable+" WHERE "+model.UserTable+".id=$1 LIMIT 1", id)
	if err != nil && err != sql.ErrNoRows {
		u.logger.Error(err)
		return err, nil
	}
	if rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Description,
			&user.Avatar,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			u.logger.Error(err)
			return err, nil
		}
	}
	err, role := u.getUserRole(user.ID)
	if err != nil && err != sql.ErrNoRows {
		u.logger.Error(err)
		return err, nil
	}
	user.Role = model.Role{ID: role.ID, Title: role.Title, Description: role.Description}
	return nil, &user
}

func (u UserPostgres) getUserRole(userId uint64) (error, *model.Role) {
	var role model.Role
	rows, err := u.db.Queryx("SELECT roles.id, roles.title, roles.description FROM "+model.RoleTable+
		" inner join "+model.UserRoleTable+" ur on roles.id = ur.role_id"+
		" where ur.user_id = $1 limit 1", userId)
	if err != nil && err != sql.ErrNoRows {
		u.logger.Error(err)
		return err, nil
	}
	if rows.Next() {
		err := rows.Scan(
			&role.ID,
			&role.Title,
			&role.Description,
		)
		if err != nil {
			u.logger.Error(err)
			return err, nil
		}
	}
	return nil, &role
}
