package main

import (
	"./config"
	"./controllers"
	"os"

	"github.com/kataras/iris"
)

const PathWeb = "./web/dist/"

func main() {
	db, session, app := config.Init()
	defer db.Close()
	defer session.Close()

	// SPA (git submodule, dist folder)
	app.StaticWeb("/", PathWeb)

	// Route -> /api/*
	api := app.Party("/api")

	// Auth
	controllers.AuthRoute(api)

	app.Run(iris.Addr(os.Getenv(":" + "PORT")))
}
