package controllers

import (
	"../config"
	"../models"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

func GroupsRoute(router router.Party) {
	// Route -> /api/groups/*
	groupsRoute := router.Party("/groups")

	groupsRoute.Get("/", handleGet)
	groupsRoute.Post("/", handlePost)
}

func handleGet(ctx context.Context) {
	db := config.GetDb()

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

	db.
		Select("DISTINCT groups.*").
		Preload("Creator").
		Preload("Users").
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("DISTINCT tasks.*").
				Group("tasks.group_id").
				Order("tasks.updated_at desc").
				Preload("User")
		}).
		Joins("LEFT JOIN group_user gu ON gu.group_id = groups.id").
		Where("groups.creator_id = ? OR gu.user_id = ?", userId, userId).
		Order("groups.updated_at desc").
		Limit(count).
		Offset(offset).
		Find(&groups)

	ctx.JSON(iris.Map{"result": "Success", "groups": groups})
}

func handlePost(ctx context.Context) {
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
		ctx.JSON(iris.Map{"result": isBlank, "group": group})
	} else {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Error"})
	}
}
