package app

import (
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"go-instaloader/WebServer/route"
	"go-instaloader/config"
	"go-instaloader/utils/rlog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var app *iris.Application

func handleSignal(server net.Listener) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		rlog.Infof("got signal [%s], exiting now", s)
		if err := server.Close(); nil != err {
			rlog.Error("server close failed", err.Error())
		}

		rlog.Info("exited")
		os.Exit(0)
	}()
}

func IrisLogFunc(ctx *context.Context, latency time.Duration) {
	var ip, method, path string

	status := ctx.GetStatusCode()
	method = ctx.Method()
	path = ctx.Request().URL.RequestURI()
	if method == "OPTIONS" {
		return
	}

	line := fmt.Sprintf("%4v %s %s %v %s", latency, ip, method, status, path)
	if context.StatusCodeNotSuccessful(status) {
		body, _ := ctx.GetBody()
		rlog.Error(line, string(body))
		return
	}
}

func IrisInit() {
	app = iris.New()
	//recover
	app.Use(recover.New())
	//logger
	app.Logger().SetLevel("info")

	app.Logger().SetOutput(rlog.GetLogger().GetWriter())
	irisLogConfig := logger.DefaultConfig()
	irisLogConfig.LogFuncCtx = IrisLogFunc
	app.Use(logger.New(irisLogConfig))
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           600,
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost /*, iris.MethodOptions, iris.MethodHead, iris.MethodDelete, iris.MethodPut*/},
		AllowedHeaders:   []string{"*"},
	}))

	app.AllowMethods(iris.MethodOptions)

	app.Use(func(c *context.Context) {
		c.Next()
	})

	route.RegisterRoutes(app)
}

func IrisStart() {
	// start web service
	listener, err := net.Listen("tcp4", config.Instance.Host)
	if err != nil {
		rlog.Error(err)
		os.Exit(1)
	}

	handleSignal(listener)
	if err := app.Run(iris.Listener(listener), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		EnableOptimizations:               true,
		TimeFormat:                        "2006-01-02 15:04:05",
		Charset:                           "UTF-8",
	})); err != nil {
		rlog.Error(err)
		os.Exit(1)
	}
}
