package main

import (
	"./models"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

var db *gorm.DB

const PathWeb = "./web/dist/"

func main() {
	db, app := Init()
	defer db.Close()

	// SPA (git submodule, dist folder)
	app.StaticWeb("/", PathWeb)

	routeApi(app)

	app.Run(iris.Addr(":8080"))
}

func routeApi(app *iris.Application) {
	api := app.Party("/api")

	api.Get("/users", func(ctx context.Context) {
		ctx.Text("Success")
	})

	// FIXME Example
	//api.Get("/api/users", func(ctx context.Context) {
	//	var user models.User
	//	db.First(&user, 1)
	//
	//	ctx.JSON(user)
	//})
}

func Init() (*gorm.DB, *iris.Application) {
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

	/* - Initialization Iris - */
	app := iris.New()
	app.Logger().SetLevel(os.Getenv("LOGGER"))

	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// Order of those calls doesn't matter, `UseGlobal` and `DoneGlobal`
	// are applied to existing routes and future routes.
	app.UseGlobal(beforeRoute)

	return db, app
}

func beforeRoute(ctx iris.Context) {
	var user models.User
	db.First(&user, 1)
	println(user.ID)

	if user.ID < 1 {
		return
	}

	ctx.Values().Set("user", user)
	println("Before")
	ctx.Next()
}
