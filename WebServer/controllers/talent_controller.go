package controllers

import (
	"github.com/kataras/iris/v12"
	"go-instaloader/WebServer/caches"
	"go-instaloader/WebServer/services"
	"go-instaloader/models"
	"go-instaloader/models/request"
	"go-instaloader/models/response"
)

var TalentController = new(talentController)

type talentController struct{}

func (t *talentController) GetTalent(ctx iris.Context) {
	req := request.GetBodyToMap(ctx)
	username := request.GetValueString(req, "username")

	if username == "" {
		response.FailWithMessageV2("username can't be empty", ctx)
		return
	}

	talent := caches.TalentCache.Get(username)
	if talent == nil {
		response.FailWithMessageV2("talent not found", ctx)
		return
	}

	response.OkWithMessageV2("ok", talent, ctx)
}

func (t *talentController) UpdateTalent(ctx iris.Context) {
	var talent *models.Talent
	if err := ctx.ReadBody(&talent); err != nil {
		response.FailWithMessageV2("failed to parse request body", ctx)
		return
	}

	if err := services.TalentService.UpdateTalentData(talent); err != nil {
		response.FailWithMessageV2("failed to update talent", ctx)
		return
	}
	response.OkWithMessageV2("ok", nil, ctx)
}

func (t *talentController) DeleteTalent(ctx iris.Context) {
	req := request.GetBodyToMap(ctx)
	username := request.GetValueString(req, "username")

	if username == "" {
		response.FailWithMessageV2("failed to get username", ctx)
		return
	}

	talent := caches.TalentCache.Get(username)
	if talent == nil {
		response.FailWithMessageV2("talent not found", ctx)
		return
	}

	if err := services.TalentService.DeleteTalentData(talent); err != nil {
		response.FailWithMessageV2("failed to delete talent", ctx)
		return
	}
	response.OkWithMessageV2("ok", nil, ctx)
}
