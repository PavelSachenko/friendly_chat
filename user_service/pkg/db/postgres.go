package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pavel/user_service/config"
)

type DB struct {
	*sqlx.DB
}

func InitPostgres(cfg *config.Config) (error, *DB) {
	connection, err := sqlx.Connect(cfg.DB.Driver, fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.Driver,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
		cfg.DB.SSLMode),
	)
	db := &DB{
		connection,
	}

	if err != nil {
		return err, nil
	}

	err = db.Ping()
	if err != nil {
		return err, nil
	}

	return nil, db
}
