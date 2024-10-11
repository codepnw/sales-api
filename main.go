package main

import (
	"github.com/codepnw/sales-api/database"
	"github.com/joho/godotenv"
)

const envFile string = "dev.env"

func main() {
	if err := godotenv.Load(envFile); err != nil {
		panic("failed load env file")
	}

	if err := database.NewConnect(); err != nil {
		panic(err)
	}
}