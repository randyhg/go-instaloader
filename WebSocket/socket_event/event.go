package socket_event

import (
	socketio "github.com/googollee/go-socket.io"
	"go-instaloader/models/constants"
	"go-instaloader/utils/rlog"
)

func InitSocketEvent(server *socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		rlog.Info("/:connected=====================================", s.ID())
		return nil
	})

	server.OnEvent("/", constants.TestEvent, SocketEvent.Test)
	server.OnEvent("/", constants.FetchEvent, SocketEvent.FetchData)
	server.OnEvent("/", constants.VerifyEvent, SocketEvent.VerifyData)

	server.OnError("/", func(s socketio.Conn, e error) {
		rlog.Error("/:error ", e)
		//s.LeaveAll()
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		rlog.Info("/:closed=====================================", s.ID(), reason)
		s.LeaveAll()
	})
}
