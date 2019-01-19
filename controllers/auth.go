package controllers

import (
	"../models"
	"../utils"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

func AuthRoute(router router.Party) {
	// Route -> /api/groups/*
	authRoute := router.Party("/auth")

	authRoute.Post("/", handleLogin)
	authRoute.Post("/logout", handleLogout)
	authRoute.Post("/logout/all", handleLogoutAll)
	authRoute.Get("/tokens", handleTokens)
	authRoute.Get("/me", handleMe)
}

func handleMe(ctx context.Context) {
	ctx.JSON(iris.Map{"result": "Success", "user": models.GetCurrentUser()})
}

func handleLogin(ctx context.Context) {
	login, password := ctx.PostValue("login"), ctx.PostValue("password")

	if login == "" || password == "" {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Login or password is empty"})
		return
	}

	user := models.GetUserByLogin(login)

	// TODO Protect Brute-force

	if user.ID <= 0 || !utils.CheckPasswordHash(password, user.Password) {
		ctx.StatusCode(422)
		ctx.JSON(iris.Map{"error": "Login or password is incorrect"})
		return
	}

	token, err := user.AddToken(ctx.RemoteAddr())

	if err == nil {
		ctx.JSON(iris.Map{"result": "Success", "token": token})
	} else {
		ctx.JSON(iris.Map{"error": "Error"})
	}
}

func handleLogout(ctx context.Context) {
	user := models.GetCurrentUser()
	token := ctx.GetHeader("Auth")
	err := user.RemoveToken(token)

	if err != nil {
		ctx.StatusCode(422)
		ctx.JSON(iris.Map{"error": "Error"})
		return
	}

	ctx.JSON(iris.Map{"result": "Success"})
}

func handleLogoutAll(ctx context.Context) {
	user := models.GetCurrentUser()
	err := user.RemoveTokens()

	if err != nil {
		ctx.StatusCode(422)
		ctx.JSON(iris.Map{"error": "Error"})
		return
	}

	ctx.JSON(iris.Map{"result": "Success"})
}

func handleTokens(ctx context.Context) {
	tokens, err := models.GetCurrentUser().GetTokens()

	if err == nil {
		ctx.JSON(iris.Map{"result": "Success", "tokens": tokens})
	} else {
		ctx.JSON(iris.Map{"error": "Error"})
	}
}
