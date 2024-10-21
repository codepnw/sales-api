package main

import (
	"github.com/codepnw/sales-api/config"
	"github.com/codepnw/sales-api/database"
	"github.com/codepnw/sales-api/routes"
	"github.com/gin-gonic/gin"
)

const (
	configPath string = "."
	configFile string = "app_config"
)

func main() {
	config.InitTimezone()
	cfg := config.InitConfig(configPath, configFile)

	if err := database.NewPostgresConnect(cfg); err != nil {
		panic(err)
	}

	app := gin.Default()
	routes.Setup(app, cfg.App().Version())

	app.Run(cfg.App().Port())
}
