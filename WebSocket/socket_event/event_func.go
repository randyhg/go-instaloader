package socket_event

import (
	"context"
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"go-instaloader/WebServer/services"
	"go-instaloader/WebSocket/socket_resp"
	"go-instaloader/models/request"
	"go-instaloader/utils/rlog"
)

var SocketEvent = new(socketEvent)

type socketEvent struct{}

func (s *socketEvent) FetchData(sck socketio.Conn, msg string) {
	fmt.Printf("%+v\n", msg)
	rlog.Info(sck.ID())
	var req request.FetchRequest
	if err := json.Unmarshal([]byte(msg), &req); err != nil {
		rlog.Error(err)
		socket_resp.DoEmitError(sck, err)
		return
	}

	socket_resp.DoEmit(sck, "fetch", "fetch talents on progress")

	if err := services.FetchService.FetchTalent(req.FetchRange, context.Background(), &sck); err != nil {
		socket_resp.DoEmitError(sck, err)
		return
	}
}

func (s *socketEvent) Test(sck socketio.Conn, msg string) {
	rlog.Infof("get test: %s", msg)
	socket_resp.DoEmit(sck, "test", "pong")
}
