package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pavel/user_service/pkg/db"
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/model"
	"time"
)

type User interface {
	One(id uint64) (error, *model.SelectUser)
	All(filter model.UserFilter) (err error, users []*model.SelectUser)
	UpdateAvatar(id uint64, filename string) (error, *model.SelectUser, string)
	Update(user model.UpdateUser) (error, *model.SelectUser)
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

func (u UserPostgres) UpdateAvatar(id uint64, filename string) (error, *model.SelectUser, string) {
	rows, err := u.db.Queryx(fmt.Sprintf("SELECT avatar FROM %s WHERE id = $1", model.UserTable), id)
	var oldAvatar string
	if rows.Next() {
		err = rows.Scan(&oldAvatar)
		if err != nil {
			u.logger.Error(err)
			return err, nil, ""
		}
	}
	rows.Close()
	rows, err = u.db.Queryx(fmt.Sprintf("UPDATE %s SET avatar = $1 WHERE id = $2 RETURNING id,username,description,avatar,created_at", model.UserTable), filename, id)
	if err != nil {
		u.logger.Error(err)
		return err, nil, ""
	}
	err, user := u.scanUserStruct(rows)

	return err, user, oldAvatar
}

func (u UserPostgres) Update(userUpdate model.UpdateUser) (error, *model.SelectUser) {
	row, err := u.db.Queryx(
		fmt.Sprintf("UPDATE %s SET description = $1, updated_at = $2 WHERE id = $3 RETURNING id,username,description,avatar,created_at", model.UserTable),
		userUpdate.Description,
		time.Now(),
		userUpdate.ID,
	)
	if err != nil {
		u.logger.Error(err)
		return err, nil
	}

	return u.scanUserStruct(row)
}

func (u UserPostgres) scanUserStruct(row *sqlx.Rows) (error, *model.SelectUser) {
	var user model.SelectUser
	if row.Next() {
		err := row.StructScan(&user)
		if err != nil {
			u.logger.Error(err)
			return err, nil
		}
	}
	return nil, &user
}
