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

	groupsRoute.Get("/", handleGroupsGet)
	groupsRoute.Post("/", handleGroupPost)
	groupsRoute.Delete("/{id:int}", handleGroupDelete)
}

func handleGroupsGet(ctx context.Context) {
	offset, err := ctx.URLParamInt("offset")
	if err != nil || offset < 0 {
		offset = 0
	}

	count, err := ctx.URLParamInt("count")
	if err != nil || count > 100 || count < 0 {
		count = 30
	}

	var groups []models.Group
	userId := models.GetCurrentUser().ID

	models.GetAllGroups(&groups, userId, count, offset)

	ctx.JSON(iris.Map{"result": "Groups received", "groups": groups})
}

func handleGroupPost(ctx context.Context) {
	name, description := ctx.PostValue("name"), ctx.PostValue("description")

	if name == "" {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Name is required"})
		return
	}

	group := models.Group{
		Name: name,
		Description: description,
		CreatorID: models.GetCurrentUser().ID,
	}

	db := config.GetDb()
	isBlank := db.NewRecord(group)

	if isBlank {
		db.Create(&group)
		db.Model(&group).Association("Users").Append(models.GetCurrentUser())
		group.Creator = models.GetCurrentUser()
		ctx.JSON(iris.Map{"result": "Group created", "group": group})
	} else {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Create error"})
	}
}

func handleGroupDelete(ctx context.Context) {
	groupId := ctx.Params().Get("id")

	var group models.Group
	db := config.GetDb()
	db.First(&group, groupId)

	if group.ID < 1 {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Group does not exist"})
		return
	}

	currentUser := models.GetCurrentUser()

	if group.CreatorID != currentUser.ID && !currentUser.IsAdmin {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "No permission to delete"})
		return
	}

	db.Delete(&group)
	ctx.JSON(iris.Map{"result": "Group deleted"})
}
