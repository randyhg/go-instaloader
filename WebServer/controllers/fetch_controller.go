package controllers

import (
	"go-instaloader/WebServer/services"
	"go-instaloader/models/request"
	"go-instaloader/models/response"

	"github.com/kataras/iris/v12"
)

var FetchController = new(fetchController)

type fetchController struct {
}

func (f *fetchController) FetchTalentData(ctx iris.Context) {
	req := request.GetBodyToMap(ctx)
	fetchRange := request.GetValueString(req, "fetch_range")

	if err := services.FetchService.FetchTalent(fetchRange, ctx); err != nil {
		response.FailWithMessageV2(err.Error(), ctx)
		return
	}

	response.OkWithMessageV2("fetch talents on progress", nil, ctx)
}
