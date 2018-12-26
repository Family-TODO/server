package main

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

const PathWeb = "./web/dist/"

func main() {
	/* - Import Environment - */
	err := godotenv.Load()
	if err != nil {
		panic("env not loaded")
	}

	/* - Connect to Database - */
	db, err := gorm.Open("sqlite3", os.Getenv("DATABASE_PATH"))
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	/* - Initialization Iris - */
	app := iris.New()
	app.Logger().SetLevel(os.Getenv("LOGGER"))

	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// SPA (git submodule, dist folder)
	app.StaticWeb("/", PathWeb)

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
