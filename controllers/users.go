package controllers

import (
	"../models"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

func UsersRoute(router router.Party) {
	// Router -> /api/users/*
	usersRoute := router.Party("/users")

	usersRoute.Get("/", handleGet)
}

func handleGet(ctx context.Context) {
	users := models.GetUsers()

	ctx.JSON(iris.Map{"result": "Users received", "users": users})
}
