package main

import (
	"./config"
	"./models"
)

func main() {
	db, session, _ := config.Init()
	defer db.Close()
	defer session.Close()

	// Migrate the schema
	models.Migrate(db)
}
