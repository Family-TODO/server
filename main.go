package main

import (
	"os"

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

	/* - Initialization Iris - */
	app := iris.New()
	app.Logger().SetLevel(os.Getenv("LOGGER"))

	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	app.StaticWeb("/", PathWeb)



	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
