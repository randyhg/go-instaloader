package socket_app

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"go-instaloader/WebServer/app"
	"go-instaloader/WebSocket/socket_event"
	"go-instaloader/utils/rlog"
	"log"
)

var Server = socketio.NewServer(nil)

func SocketStart() {
	socket := iris.New()
	socket.Logger().SetLevel("info")

	socket.Logger().SetOutput(rlog.GetLogger().GetWriter())
	irisLogConfig := logger.DefaultConfig()
	irisLogConfig.LogFuncCtx = app.IrisLogFunc
	socket.Use(logger.New(irisLogConfig))

	socket_event.InitSocketEvent(Server)

	go func() {
		if err := Server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer Server.Close()

	socket.HandleMany("GET POST", "/socket.io/{any:path}", iris.FromStd(Server))
	if err := socket.Run(
		iris.Addr(":8091"),
		iris.WithoutPathCorrection,
		iris.WithoutServerError(iris.ErrServerClosed),
	); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
