package controllers

import (
	"../config"
	"../models"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

func AuthRoute(router router.Party) {
	router.Post("/auth", handleLogin)
}

func handleLogin(ctx context.Context) {
	db := config.GetDB()
	login, password := ctx.PostValue("login"), ctx.PostValue("password")

	if login == "" || password == "" {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Login or password is empty"})
		return
	}

	var user models.User
	db.Where("login = ?", login).First(&user)

	if user.ID <= 0 || !models.CheckPasswordHash(password, user.Password) {
		ctx.StatusCode(422)
		ctx.JSON(iris.Map{"error": "Login or password is incorrect"})
		return
	}

	session := config.GetSession()
	sess := session.Start(ctx)
	sess.Set("isAuth", true)
	sess.Set("ID", user.ID)
	ctx.JSON(iris.Map{"result": "Success"})
}
