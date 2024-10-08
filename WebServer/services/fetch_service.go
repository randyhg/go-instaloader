package services

import (
	"errors"
	"github.com/kataras/iris/v12"
	"go-instaloader/utils/rlog"
)

var FetchService = new(fetchService)

type fetchService struct{}

func (s *fetchService) FetchTalent(fetchRange string, ctx iris.Context) error {
	sheet := newSheetService()
	if sheet == nil {
		return errors.New("unable to get talents")
	}
	go func() {
		talents, err := sheet.GetTalents(ctx, fetchRange)
		if err != nil {
			rlog.Errorf("unable to get talents: %s", err.Error())
			ErrorHandler(err)
		} else {
			if talents == nil {
				rlog.Info("no talents found")
				ErrorHandler(errors.New("no talents found"))
				return
			}

			rlog.Info("talents pushed to redis")
		}
	}()
	return nil
}
