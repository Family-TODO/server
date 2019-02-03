package main

import (
	"./config"
	"./models"
)

func main() {
	config.NewConfig()

	// Migrate the schema
	models.Migrate(config.Db)
}
