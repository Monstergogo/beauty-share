package server

import (
	"context"
	"fmt"
	pb "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/internal/app"
	"github.com/Monstergogo/beauty-share/internal/repo_interface"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServerRouter struct {
	Path    string
	Method  util.HttpMethod
	Handler func(context *gin.Context)
}

type MicroServer struct {
	GinServer *gin.Engine
}

func (m MicroServer) RegisterGinRouter(router HttpServerRouter) {
	switch router.Method {
	case util.HttpMethodGet:
		m.GinServer.GET(router.Path, router.Handler)
	case util.HttpMethodPost:
		m.GinServer.POST(router.Path, router.Handler)
	case util.HttpMethodPut:
		m.GinServer.PUT(router.Path, router.Handler)
	case util.HttpMethodDelete:
		m.GinServer.DELETE(router.Path, router.Handler)
	case util.HttpMethodPatch:
		m.GinServer.PATCH(router.Path, router.Handler)
	}
}

func InitServer() MicroServer {
	logger.InitLogger(logger.LogConf{
		LogFilepath: util.LogPath,
		ErrFilepath: util.ErrPath,
	})
	ginServer := gin.Default()
	db.InitMongoDB()
	return MicroServer{GinServer: ginServer}
}

func (m MicroServer) RunServer() {
	httpServer := http.Server{
		Addr:    ":5008",
		Handler: m.GinServer,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Error("http server listen err", zap.Any("err", err))
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 5018))
		if err != nil {
			logger.GetLogger().Error("grpc server listen err", zap.Any("err", err))
		}
		srv := grpc.NewServer()
		pb.RegisterShareServiceServer(srv, &app.ShareServiceImpl{
			MongoRepo: repo_interface.MongoRepoProvider(),
		})
		if err = srv.Serve(listener); err != nil {
			logger.GetLogger().Error("grpc server serve err", zap.Any("err", err))
		}
	}()
	waitSignalClose(httpServer)
}

func waitSignalClose(server http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	db.DisconnectMongoDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.GetLogger().Error("grace shutdown err", zap.Any("err", err))
	}
	logger.GetLogger().Info("server shutdown success")
}
