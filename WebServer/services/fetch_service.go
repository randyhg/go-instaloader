package services

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"go-instaloader/models"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/rlog"
)

var FetchService = new(fetchService)

type fetchService struct{}

func (s *fetchService) FetchTalent(fetchRange string, ctx iris.Context) error {
	talents, err := SheetService.GetTalents(ctx, fetchRange)
	if err != nil {
		rlog.Error("unable to get talents:", err.Error())
		return err
	}

	if talents == nil {
		rlog.Info("no talents found")
		return nil
	}

	for _, talent := range talents {
		if len(talent.Uuid) > 0 {
			byt, err := json.Marshal(&talent)
			if err != nil {
				rlog.Error(err)
			} else {
				fwRedis.RedisQueue().LPush(ctx, models.RedisJobQueueKey, string(byt))
			}
		}
	}
	rlog.Info("talents data pushed to redis")
	return nil
}
