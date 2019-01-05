package config

import (
	"os"

	"github.com/dgraph-io/badger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

var (
	db       *gorm.DB
	badgerDb *badger.DB
)

func Init() (*gorm.DB, *badger.DB, *iris.Application) {
	/* - Import Environment - */
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	/* - Connect to Database - */
	db, err = gorm.Open("sqlite3", os.Getenv("DATABASE_PATH"))

	if err != nil {
		panic(err)
	}

	/* - Key-value Database - */
	opts := badger.DefaultOptions
	opts.Dir = os.Getenv("BADGER_PATH")
	opts.ValueDir = os.Getenv("BADGER_PATH")
	badgerDb, err = badger.Open(opts)
	if err != nil {
		panic(err)
	}

	// Close and unlock the database when control+C/cmd+C pressed
	iris.RegisterOnInterrupt(func() {
		db.Close()
	})

	/* - Initialization Iris - */
	app := iris.New()
	app.Logger().SetLevel(os.Getenv("LOGGER"))

	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	return db, badgerDb, app
}

func GetDb() *gorm.DB {
	return db
}

func GetBadgerDb() *badger.DB {
	return badgerDb
}
