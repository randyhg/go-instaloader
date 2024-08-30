package services

import (
	"context"
	"errors"
	socketio "github.com/googollee/go-socket.io"
	"go-instaloader/WebSocket/socket_resp"
	"go-instaloader/models/constants"
	"go-instaloader/utils/rlog"
)

var FetchService = new(fetchService)

type fetchService struct{}

func (s *fetchService) FetchTalent(fetchRange string, ctx context.Context, sck *socketio.Conn) error {
	sheet := newSheetService()
	if sheet == nil {
		return errors.New("unable to get talents")
	}
	go func() {
		talents, err := newSheetService().GetTalents(ctx, fetchRange)
		if err != nil {
			rlog.Errorf("unable to get talents: %s", err.Error())
			//return err
		}

		if talents == nil {
			rlog.Info("no talents found")
			//return nil
		}

		rlog.Info("talents pushed to redis")
		socket_resp.DoEmit(*sck, constants.FetchEvent, "fetch talents success")
	}()
	return nil
}
