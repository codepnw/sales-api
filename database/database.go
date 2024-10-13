package database

import (
	"fmt"
	"os"

	"github.com/codepnw/sales-api/entities"
	"github.com/codepnw/sales-api/pkg/logs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func NewConnect() error {
	DBHost := os.Getenv("MYSQL_HOST")
	DBPort := os.Getenv("MYSQL_PORT")
	DBUser := os.Getenv("MYSQL_USER")
	DBPassword := os.Getenv("MYSQL_PASSWORD")
	DBName := os.Getenv("MYSQL_DBNAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPassword, DBHost, DBPort, DBName,
	)
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("failed connection database")
	}

	db = connection
	logs.Info("database connected successfully")

	if err := AutoMigrate(connection); err != nil {
		return err
	}
	return nil
}

func AutoMigrate(connection *gorm.DB) error {
	if err := connection.Debug().AutoMigrate(
		&entities.Cashier{},
		&entities.Category{},
		&entities.Payment{},
		&entities.PaymentType{},
		&entities.Product{},
		&entities.Discount{},
		&entities.Order{},
	); err != nil {
		return err
	}
	return nil
}

func GetDB() *gorm.DB {
	return db
}
