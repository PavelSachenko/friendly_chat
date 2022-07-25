package repository

import (
	"database/sql"
	"fmt"
	"github.com/pavel/user_service/pkg/db"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/model"
)

type User interface {
	One(id uint64) (error, *model.SelectUser)
	All(filter model.UserFilter) (err error, users []*model.SelectUser)
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

func (u UserPostgres) All(filter model.UserFilter) (err error, users []*model.SelectUser) {
	rawSql := fmt.Sprintf("SELECT id, username, description, avatar, created_at FROM %s WHERE username like $1 AND id <> $2 LIMIT %d OFFSET %d", model.UserTable, filter.Limit, filter.Offset)
	rows, err := u.db.Queryx(rawSql, filter.Username+"%", filter.OwnerUserId)
	if err != nil {
		u.logger.Warnf("Bad sql request: err %v", err)
		return err, nil
	}
	for rows.Next() {
		var user model.SelectUser
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Description,
			&user.Avatar,
			&user.CreatedAt,
		)
		if err != nil {
			u.logger.Warnf("Can't scan struct: err %v", err)
			return err, nil
		}
		users = append(users, &user)
	}
	rows.Close()

	return nil, users
}

func (u UserPostgres) One(id uint64) (error, *model.SelectUser) {
	var user model.SelectUser
	rows, err := u.db.Queryx("SELECT id,username,description,avatar,created_at FROM "+model.UserTable+" WHERE "+model.UserTable+".id=$1 LIMIT 1", id)
	if err != nil && err != sql.ErrNoRows {
		u.logger.Error(err)
		return err, nil
	}
	if rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Description,
			&user.Avatar,
			&user.CreatedAt,
		)
		if err != nil {
			u.logger.Error(err)
			return err, nil
		}
	}
	rows.Close()
	return nil, &user
}
