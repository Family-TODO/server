package main

import (
	"./config"
	"./controllers"
	"./models"

	"os"

	"github.com/kataras/iris"
)

const PathWeb = "./web/dist/"

var (
	allowNotAuthRoutesName = [2]string{
		"GET/*file",
		"POST/api/auth",
	}
	blockAuthRoutesName = [2]string{
		"POST/api/auth",
	}
)

func main() {
	db, badgerDb, app := config.Init()
	defer db.Close()
	defer badgerDb.Close()

	// SPA (git submodule, dist folder)
	app.StaticWeb("/", PathWeb)

	// Order of those calls doesn't matter, `UseGlobal` and `DoneGlobal`
	// are applied to existing routes and future routes.
	app.UseGlobal(beforeRoute)

	// Route -> /api/*
	api := app.Party("/api")

	// Routes
	controllers.AuthRoute(api)
	controllers.GroupsRoute(api)
	controllers.UsersRoute(api)

	// Run server
	if os.Getenv("APP_MODE") == "release" || os.Getenv("TLS_ENABLE") == "1" || os.Getenv("TLS_ENABLE") == "true" {
		if os.Getenv("TLS_AUTO") == "1" || os.Getenv("TLS_AUTO") == "true" {
			app.Run(iris.AutoTLS(os.Getenv("TLS_ADDR"), os.Getenv("TLS_DOMAIN"), os.Getenv("TLS_EMAIL")))
		} else {
			app.Run(iris.TLS(os.Getenv("TLS_ADDR"), os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY")))
		}
	} else {
		app.Run(iris.Addr(os.Getenv("APP_ADDR")))
	}
}

// Guard
func beforeRoute(ctx iris.Context) {
	// Check header token
	authTokenHeader := ctx.GetHeader("Auth")
	isAuth := models.ValidateUserToken(authTokenHeader)

	currentRouteName := ctx.GetCurrentRoute().Name()

	if isAuth && existRouteName(currentRouteName, blockAuthRoutesName) {
		ctx.StatusCode(iris.StatusMethodNotAllowed)
		ctx.JSON(map[string]string{"error": "Access denied for authorized users"})
		return
	}

	if !isAuth && !existRouteName(currentRouteName, allowNotAuthRoutesName) {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "Auth is required"})
		return
	}

	ctx.Next()
}

func existRouteName(currentRouteName string, routesName [2]string) bool {
	for _, routeName := range routesName {
		if currentRouteName == routeName {
			return true
		}
	}

	return false
}
