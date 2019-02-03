package config

import (
	"../telegram"
	"../utils"

	"os"
	"strconv"

	"github.com/dgraph-io/badger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

var (
	BadgerDb *badger.DB
	TlgBot   *telegram.TlgBot
	Db       *gorm.DB
)

func NewConfig() *iris.Application {
	/* - Init Environment .env - */
	initEnvironment()

	/* - TelegramBot for notification - */
	TlgBot = initTelegram()

	/* - Connect to sqlite3 Database - */
	Db = initDatabase()

	/* - Connect to Badger Database - */
	BadgerDb = initBadgerDatabase()

	/* - Init Iris Server - */
	app := initIrisServer()

	/* - Close everything on close server - */
	iris.RegisterOnInterrupt(func() {
		Db.Close()
		BadgerDb.Close()

		if utils.EnvIsRelease() {
			_, _ = TlgBot.SendMessage("Server Stopped")
		}
	})

	return app
}

// Import config from .env file
func initEnvironment() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

// Get config from .env file and
// set init
func initTelegram() *telegram.TlgBot {
	u, err := strconv.ParseUint(os.Getenv("TLG_USER_ID"), 10, 64)
	if err != nil {
		panic(err)
	}

	return telegram.NewTelegram(os.Getenv("TLG_TOKEN"), u)
}

// Connect to sqlite3 Database
// In debug mode - enable log mode
func initDatabase() *gorm.DB {
	db, err := gorm.Open("sqlite3", os.Getenv("DATABASE_PATH"))
	if err != nil {
		panic(err)
	}

	if !utils.EnvIsRelease() {
		db.LogMode(true)
	}

	return db
}

// Connect to Key-Value Database - badgerDB
func initBadgerDatabase() *badger.DB {
	opts := badger.DefaultOptions
	opts.Dir = os.Getenv("BADGER_PATH")
	opts.ValueDir = os.Getenv("BADGER_PATH")

	badgerDb, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}

	return badgerDb
}

// Create iris Application instance.
func initIrisServer() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel(os.Getenv("LOGGER"))

	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	return app
}
