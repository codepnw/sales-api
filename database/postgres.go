package database

import (
	"fmt"
	"os"

	"github.com/codepnw/sales-api/pkg/logs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbPostgres *sqlx.DB

func NewPostgresConnect() error {
	DBDriver := os.Getenv("PG_DRIVER")
	DBHost := os.Getenv("PG_HOST")
	DBPort := os.Getenv("PG_PORT")
	DBUser := os.Getenv("PG_USER")
	DBPassword := os.Getenv("PG_PASSWORD")
	DBName := os.Getenv("PG_DBNAME")

	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		DBDriver, DBUser, DBPassword, DBHost, DBPort, DBName,
	)

	connection, err := sqlx.Connect(DBDriver, dsn)
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
