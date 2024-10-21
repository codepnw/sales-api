package database

import (
	"fmt"

	"github.com/codepnw/sales-api/config"
	"github.com/codepnw/sales-api/pkg/logs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbPostgres *sqlx.DB

func NewPostgresConnect(cfg config.IConfig) error {
	connection, err := sqlx.Connect(cfg.DB().Driver(), cfg.DB().DSN())
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("failed connection database")
	}

	dbPostgres = connection
	logs.Info("postgres database connected successfully")

	return nil
}

func GetPostgresDB() *sqlx.DB {
	return dbPostgres
}
