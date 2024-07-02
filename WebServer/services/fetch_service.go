package services

import (
	"context"
	"encoding/json"
	"fmt"
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
		rlog.Error(fmt.Sprintf("unable to get talents: %s", err.Error()))
		return err
	}

	if talents == nil {
		rlog.Info("no talents found")
		return nil
	}

	go func() {
		for _, talent := range talents {
			if len(talent.Uuid) > 0 {
				byt, err := json.Marshal(&talent)
				if err != nil {
					rlog.Error(err)
				} else {
					fwRedis.RedisQueue().LPush(context.Background(), models.RedisJobQueueKey, string(byt))
				}
			}
		}
		rlog.Info("talents data pushed to redis")
	}()
	return nil
}
