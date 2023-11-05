package server

import (
	"context"
	"fmt"
	pb "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/init/nacos"
	"github.com/Monstergogo/beauty-share/internal/app"
	"github.com/Monstergogo/beauty-share/internal/repo_interface"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
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
	nacos.InitNacos()
	return MicroServer{GinServer: ginServer}
}

// 注册service
func registerService() error {
	namingClient := nacos.GetNacosNamingClient()
	instanceParams := vo.RegisterInstanceParam{
		Port:        util.GrpcServerPort,
		ServiceName: util.GrpcServiceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"rpc-type": "grpc", "version": "v1"},
	}
	currIp, err := util.GetCurrIp()
	if err != nil {
		return err
	}
	instanceParams.Ip = currIp
	success, err := namingClient.RegisterInstance(instanceParams)
	if err != nil {
		return err
	}
	logger.GetLogger().Info("register grpc service instance success", zap.Bool("register result", success))
	return err
}

// 注销service实例
func deRegisterService() error {
	namingClient := nacos.GetNacosNamingClient()
	instanceParams := vo.DeregisterInstanceParam{
		Port:        util.GrpcServerPort,
		ServiceName: util.GrpcServiceName,
		Ephemeral:   true,
	}
	currIp, err := util.GetCurrIp()
	if err != nil {
		return err
	}
	instanceParams.Ip = currIp
	success, err := namingClient.DeregisterInstance(instanceParams)
	if err != nil {
		return err
	}
	logger.GetLogger().Info("deregister grpc service success", zap.Any("deregister result", success))
	return err
}

func (m MicroServer) RunServer() {
	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", util.GinServerPort),
		Handler: m.GinServer,
	}

	serverStartErrChan := make(chan interface{}, 2)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Error("http server listen err", zap.Any("err", err))
			serverStartErrChan <- err
			return
		}
	}()

	srv := grpc.NewServer()
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", util.GrpcServerPort))
		if err != nil {
			logger.GetLogger().Error("grpc server listen err", zap.Any("err", err))
			return
		}
		pb.RegisterShareServiceServer(srv, &app.ShareServiceImpl{
			MongoRepo: repo_interface.MongoRepoProvider(),
		})
		if err = srv.Serve(listener); err != nil {
			logger.GetLogger().Error("grpc server serve err", zap.Any("err", err))
			serverStartErrChan <- err
			return
		}
	}()
	time.Sleep(1 * time.Second)
	// 注册grpc服务到nacos
	if err := registerService(); err != nil {
		logger.GetLogger().Error("register service instance err", zap.Any("err_msg", err))
		serverStartErrChan <- err
		return
	}
	waitSignalClose(&httpServer, srv, serverStartErrChan)
}

func waitSignalClose(ginServer *http.Server, grpcServer *grpc.Server, serverStartErrChan chan interface{}) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		logger.GetLogger().Info("system quit")
	case errInfo := <-serverStartErrChan:
		panic(errInfo)
	}

	db.DisconnectMongoDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := ginServer.Shutdown(ctx); err != nil {
			logger.GetLogger().Error("gin server grace shutdown err", zap.Any("err", err))
		}
		logger.GetLogger().Info("gin shutdown")
		return
	}()

	go func() {
		defer wg.Done()
		grpcServer.GracefulStop()
		deRegisterService()
		return
	}()
	wg.Wait()
	logger.GetLogger().Info("server shutdown success")
}
