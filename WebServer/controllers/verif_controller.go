package controllers

import (
	"github.com/kataras/iris/v12"
	"go-instaloader/WebServer/services"
	"go-instaloader/models"
	"go-instaloader/models/request"
	"go-instaloader/models/response"
)

var VerifController = new(verifController)

type verifController struct {
}

func (c *verifController) VerifProfileAndStoryTalents(ctx iris.Context) {
	req := request.GetBodyToMap(ctx)
	storyLimit := request.GetValueIntDefault(req, "story_limit", models.DefaultStoryLimit)
	url := request.GetValueString(req, "url")

	if err := services.VerifService.VerifTalentService(storyLimit, url, ctx); err != nil {
		response.FailWithMessageV2(err.Error(), ctx)
		return
	}
	response.OkWithMessageV2("the job finished successfully", nil, ctx)
}
