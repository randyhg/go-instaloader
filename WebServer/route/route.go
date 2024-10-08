package route

import (
	"github.com/kataras/iris/v12"
	"go-instaloader/WebServer/controllers"
	"go-instaloader/models/response"
)

func RegisterRoutes(app *iris.Application) {
	//opts := basicauth.Options{
	//	Allow: basicauth.AllowUsers(map[string]string{
	//		config.Instance.Username: config.Instance.Password,
	//	}),
	//	Realm:        "Authorization Required",
	//	ErrorHandler: basicauth.DefaultErrorHandler,
	//}
	//auth := basicauth.New(opts)
	//app.Use(auth)

	mainGroup := app.Party("/api")
	mainGroup.Get("/ping", func(ctx iris.Context) {
		response.OkWithMessageV2("ok", "ok", ctx)
	})

	// api/fetch
	fetchGroup := mainGroup.Party("/fetch")
	{
		fetchGroup.Post("", controllers.FetchController.FetchTalentData)
	}

	// api/verif
	verifGroup := mainGroup.Party("/verification")
	{
		verifGroup.Post("", controllers.VerifController.VerifProfileAndStoryTalents)
	}

	// api/talent
	talentGroup := mainGroup.Party("/talent")
	{
		talentGroup.Get("/list", controllers.TalentController.GetTalentList)
		talentGroup.Get("/detail", controllers.TalentController.GetTalentDetail)
		talentGroup.Post("/update", controllers.TalentController.UpdateTalent)
		talentGroup.Post("/delete", controllers.TalentController.DeleteTalent)
	}
}
