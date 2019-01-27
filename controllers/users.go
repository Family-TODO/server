package controllers

import (
	"../models"

	"strconv"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

func UsersRoute(router router.Party) {
	// Router -> /api/users/*
	usersRoute := router.Party("/users")

	usersRoute.Get("/", handleGet)
	usersRoute.Put("/{id:int}", handlePut)
}

func handleGet(ctx context.Context) {
	users := models.GetUsers()

	ctx.JSON(iris.Map{"result": "Users received", "users": users})
}

func handlePut(ctx context.Context) {
	// Parse ID Param
	userId, err := strconv.ParseUint(ctx.Params().Get("id"), 10, 64)
	if err != nil || userId < 1 {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Invalid ID"})
		return
	}

	// Get the current user and the user being edited
	currentUser := models.GetCurrentUser()
	user := models.GetUserById(uint(userId))

	// No record
	if user.ID < 1 {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "User is not found"})
		return
	}

	// No rights to edit if you are not admin or are not editing yourself
	if !currentUser.IsAdmin && currentUser.ID != user.ID {
		ctx.StatusCode(iris.StatusMethodNotAllowed)
		ctx.JSON(iris.Map{"error": "No permissions"})
		return
	}

	var data = make(map[string]interface{})
	data["name"] = ctx.URLParam("name")

	isAdmin, err := ctx.URLParamBool("is_admin")
	if err == nil && currentUser.IsAdmin && currentUser.ID != user.ID {
		data["is_admin"] = isAdmin
	}

	user.Update(data)
	ctx.JSON(iris.Map{"reuslt": "Updated", "user": user})
}
