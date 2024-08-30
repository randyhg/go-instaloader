package socket_event

import (
	socketio "github.com/googollee/go-socket.io"
	"go-instaloader/utils/rlog"
)

func InitSocketEvent(server *socketio.Server) {
	rlog.Info("masuk init socket event")
	server.OnConnect("/", func(s socketio.Conn) error {
		rlog.Info("/:connected=====================================", s.ID())
		return nil
	})

	server.OnEvent("/", "test", SocketEvent.Test)

	server.OnEvent("/", "fetch", SocketEvent.FetchData)

	server.OnError("/", func(s socketio.Conn, e error) {
		rlog.Error("/:error ", e)
		s.LeaveAll()
		s.Close()
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		rlog.Info("/:closed=====================================", s.ID(), reason)
		s.LeaveAll()
	})
}
