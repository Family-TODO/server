package controllers

import (
	"../config"
	"../models"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

func GroupsRoute(router router.Party) {
	// Route -> /api/groups/*
	groupsRoute := router.Party("/groups")

	groupsRoute.Post("/", handlePost)
}

func handlePost(ctx context.Context) {
	name, description := ctx.PostValue("name"), ctx.PostValue("description")

	if name == "" {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Name is required"})
		return
	}

	group := models.Group{Name: name, Description: description, CreatorId: models.GetCurrentUser().ID}

	db := config.GetDb()
	isBlank := db.NewRecord(group)

	if isBlank {
		db.Create(&group)
		ctx.JSON(iris.Map{"result": isBlank, "group": group})
	} else {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Error"})
	}
}
