package main

import (
	"./config"
	"./controllers"
	"./models"
	"./utils"

	"os"

	"github.com/kataras/iris"
)

// Git submodule, web
const PathWeb = "./web/dist/"

// Uses for protect route
var (
	allowNotAuthRoutesName = []string{
		"GET/*file",
		"POST/api/auth",
	}
	blockAuthRoutesName = []string{
		"POST/api/auth",
	}
)

func main() {
	app := config.NewConfig()

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

	if utils.EnvIsRelease() {
		_, _ = config.TlgBot.SendMessage("Server Running")
	}

	// Run server
	startServer(app)
}

// Run iris server
// TLS server on release mode
func startServer(app *iris.Application) {
	if utils.EnvIsRelease() || utils.EnvIsTrue("TLS_ENABLE") {
		if utils.EnvIsTrue("TLS_AUTO") {
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

	// If True - models.currentUser is not empty
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

func existRouteName(currentRouteName string, routesName []string) bool {
	for _, routeName := range routesName {
		if currentRouteName == routeName {
			return true
		}
	}

	return false
}
