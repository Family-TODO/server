package config

import (
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/badger"
)

var (
	db      *gorm.DB
	session *sessions.Sessions

	publicRouteName = []string{
		"GET/*file",
		"POST/api/auth",
	}
)

func Init() (*gorm.DB, *badger.Database, *iris.Application) {
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

	/* - Sessions - */
	sessionDatabase, err := badger.New(os.Getenv("SESSION_PATH"))
	if err != nil {
		panic(err)
	}

	// Close and unlock the database when control+C/cmd+C pressed
	iris.RegisterOnInterrupt(func() {
		db.Close()
		sessionDatabase.Close()
	})

	session = sessions.New(sessions.Config{
		Cookie:       "_session",
		Expires:      30 * 24 * time.Hour, // <= 0 means unlimited life
		AllowReclaim: true,
	})

	session.UseDatabase(sessionDatabase)

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

	return db, sessionDatabase, app
}

func GetDB() *gorm.DB {
	return db
}

func GetSession() *sessions.Sessions {
	return session
}

func beforeRoute(ctx iris.Context) {
	sess := session.Start(ctx)
	isAuth, _ := sess.GetBoolean("isAuth")
	routeName := ctx.GetCurrentRoute().Name()

	if !isAuth && strings.Index(strings.Join(publicRouteName, ","), routeName) == -1 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "Auth is required"})
		return
	}

	ctx.Next()
}
