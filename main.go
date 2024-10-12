package main

import (
	"github.com/codepnw/sales-api/database"
	"github.com/codepnw/sales-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	envFile string = "dev.env"
	version string = "v1"
)

func main() {
	if err := godotenv.Load(envFile); err != nil {
		panic("failed load env file")
	}

	if err := database.NewConnect(); err != nil {
		panic(err)
	}

	app := gin.Default()
	routes.Setup(app, version)

	app.Run(":8080")
}
