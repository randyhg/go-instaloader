package socket_resp

import (
	socketio "github.com/googollee/go-socket.io"
	"go-instaloader/models/constants"
	"go-instaloader/models/response"
)

func DoEmit(s socketio.Conn, event string, msg interface{}) socketio.Conn {
	s.Emit(event, response.JsonData(msg))
	return s
}

func DoEmitError(s socketio.Conn, msg error) socketio.Conn {
	s.Emit(constants.ErrorNotice, response.JsonErrorMsg(msg.Error()))
	return s
}
