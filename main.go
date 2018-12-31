package main

import (
	"./config"
	"./controllers"

	"os"
	"strings"

	"github.com/kataras/iris"
)

const PathWeb = "./web/dist/"

var publicRouteName = []string{
	"GET/*file",
	"POST/api/auth",
}

func main() {
	db, session, app := config.Init()
	defer db.Close()
	defer session.Close()

	// SPA (git submodule, dist folder)
	app.StaticWeb("/", PathWeb)

	// Order of those calls doesn't matter, `UseGlobal` and `DoneGlobal`
	// are applied to existing routes and future routes.
	app.UseGlobal(beforeRoute)

	// Route -> /api/*
	api := app.Party("/api")

	// Auth
	controllers.AuthRoute(api)

	app.Run(iris.Addr(os.Getenv(":" + "PORT")))
}

func beforeRoute(ctx iris.Context) {
	sess := config.GetSession().Start(ctx)
	isAuth, _ := sess.GetBoolean("isAuth")
	routeName := ctx.GetCurrentRoute().Name()

	if !isAuth && strings.Index(strings.Join(publicRouteName, ","), routeName) == -1 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "Auth is required"})
		return
	}

	ctx.Next()
}
