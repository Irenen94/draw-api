package main

import (
	"aigc-image/config"
	"aigc-image/pkg/logger"
	"aigc-image/server/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//初始化logger
	logger.NewLogger(logger.WithLevel("debug"))
	logger.Logger.Info("Face Verify Gateway Start")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			logger.Logger.Info("Panic: ", zap.String("error: ", err.(string)))
		}
	}()

	api := &router.API{Router: gin.Default()}
	api.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Confs.BaseConf.HTTPPort),
		Handler:        api.Router,
		ReadTimeout:    time.Second * time.Duration(config.Confs.BaseConf.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(config.Confs.BaseConf.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}
	logger.Logger.Info("server start at ", zap.String("", time.Now().Format("2006-01-02 15:04:05")))
	e := s.ListenAndServe()

	logger.Logger.Panic("Start Server Panic: ", zap.Error(e))
	panic(e)

	//定义信号接收通道，订阅程序退出信号
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case <-c:
			logger.Logger.Info("Face Verify Gateway Exit")
			goto RETURN
		}
	}

RETURN:
	return

}
