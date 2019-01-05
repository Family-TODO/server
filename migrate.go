package main

import (
	"./config"
	"./models"
)

func main() {
	db, badgerDB, _ := config.Init()
	defer db.Close()
	defer badgerDB.Close()

	// Migrate the schema
	models.Migrate(db)
}
