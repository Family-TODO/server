package controllers

import (
	"../models"
	"../utils"

	"regexp"

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
	ctx.JSON(iris.Map{"result": "Profile received", "user": models.GetCurrentUser()})
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
		ctx.JSON(iris.Map{"result": "You are logged in", "token": token})
	} else {
		ctx.JSON(iris.Map{"error": "Authorisation Error"})
	}
}

func handleLogout(ctx context.Context) {
	user := models.GetCurrentUser()
	token := ctx.GetHeader("Auth")
	err := user.RemoveToken(token)

	if err != nil {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Delete error"})
		return
	}

	ctx.JSON(iris.Map{"result": "You are logged out"})
}

func handleLogoutAll(ctx context.Context) {
	user := models.GetCurrentUser()
	err := user.RemoveTokens()

	if err != nil {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Delete error"})
		return
	}

	ctx.JSON(iris.Map{"result": "Tokens removed"})
}

func handleTokens(ctx context.Context) {
	tokens, err := models.GetCurrentUser().GetTokens()

	for i, val := range tokens {
		match := regexp.MustCompile(`^(.+:\w{5}).+(\w{5})$`).FindStringSubmatch(val.Token)
		tokens[i].Token = match[1] + "..." + match[2]
	}

	if err == nil {
		ctx.JSON(iris.Map{"result": "Tokens received", "tokens": tokens})
	} else {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Receive error"})
	}
}
