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

	usersRoute.Get("/", handleUsersGet)
	usersRoute.Get("/{id:int}/tokens", handleTokens)
	usersRoute.Put("/{id:int}", handleUsersPut)
	usersRoute.Post("/", handleUserPost)
	//TODO Change password
	//TODO Delete
}

func handleUsersGet(ctx context.Context) {
	users := models.GetUsers()

	ctx.JSON(iris.Map{"result": "Users received", "users": users})
}

func handleUsersPut(ctx context.Context) {
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

func handleUserPost(ctx context.Context) {
	currentUser := models.GetCurrentUser()

	if !currentUser.IsAdmin {
		ctx.StatusCode(iris.StatusMethodNotAllowed)
		ctx.JSON(iris.Map{"error": "User is not found"})
		return
	}

	name, login, password := ctx.PostValue("name"), ctx.PostValue("login"), ctx.PostValue("password")
	isAdmin, err := ctx.PostValueBool("is_admin")

	if err != nil {
		isAdmin = false
	}

	if login == "" || password == "" {
		ctx.StatusCode(iris.StatusMethodNotAllowed)
		ctx.JSON(iris.Map{"error": "Login or password is empty"})
		return
	}

	if len(password) < 6 {
		ctx.StatusCode(iris.StatusMethodNotAllowed)
		ctx.JSON(iris.Map{"error": "Minimum password length 6"})
		return
	}

	user := models.User{
		Name:     name,
		Login:    login,
		IsAdmin:  isAdmin,
		Password: password,
	}

	models.CreateUser(&user)

	if user.ID < 1 {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "User not created, login is probably busy"})
		return
	}

	ctx.JSON(iris.Map{"result": "User created", "user": user})
}
