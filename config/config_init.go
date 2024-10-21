package config

import (
	"strings"
	"time"

	"github.com/codepnw/sales-api/pkg/logs"
	"github.com/spf13/viper"
)

func InitConfig(path, filename string) IConfig {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()
	// example: APP_PORT=5000 go run .
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		logs.Error(err)
		panic(err)
	}

	return &config{
		app: &app{
			port:    viper.GetString("app.port"),
			version: viper.GetString("app.version"),
		},
		db: &db{
			driver:         viper.GetString("db.driver"),
			dsn:            viper.GetString("db.dsn"),
			maxConnections: viper.GetInt("db.max_connections"),
		},
	}
}

func InitTimezone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}
