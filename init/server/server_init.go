package server

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"github.com/Monstergogo/beauty-share/conf"
	"github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/init/minio"
	"github.com/Monstergogo/beauty-share/init/nacos"
	grpc2 "github.com/Monstergogo/beauty-share/internal/app"
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
	"strings"
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
	logger.InitLogger(conf.ServerConf.Log.LogPath, conf.ServerConf.Log.ErrPath)
	db.InitMongoDB()
	minio.InitMinio()
	ginServer := gin.Default()
	return MicroServer{GinServer: ginServer}
}

// 注册service to nacos
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

// 从nacos 注销service实例
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

type registerToConsulPayload struct {
	ID                string                 `json:"ID"`
	Name              string                 `json:"Name"`
	Tags              []string               `json:"Tags"`
	Address           string                 `json:"Address"`
	Port              int                    `json:"Port"`
	Meta              map[string]string      `json:"Meta"`
	EnableTagOverride bool                   `json:"EnableTagOverride"`
	Check             map[string]interface{} `json:"Check"`
	Weights           map[string]int         `json:"Weights"`
}

// 注册grpc服务到consul
func registerServiceToConsul() error {
	payload := registerToConsulPayload{
		Name: util.GrpcServiceName,
		Port: util.GrpcServerPort,
		Tags: []string{"share", "v1"},
		Meta: map[string]string{"version": "0.1.1"},
	}

	currIp, err := util.GetCurrIp()
	if err != nil {
		return err
	}
	payload.Address = currIp
	// 配置health check http接口
	healthCheck := map[string]interface{}{
		"DeregisterCriticalServiceAfter": "90m",
		"HTTP":                           fmt.Sprintf("http://%s:%d/v1/ping", currIp, util.GinServerPort),
		"Interval":                       "10s",
		"Timeout":                        "5s",
	}
	payload.Check = healthCheck
	registerUrl := fmt.Sprintf("%s/v1/agent/service/register?replace-existing-checks=true", conf.ServerConf.Consul.Endpoint)
	payloadMarshal, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("PUT", registerUrl, strings.NewReader(string(payloadMarshal)))
	req.Header.Add("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	return err
}

// 从consul取消register
func deregisterServiceToConsul() error {
	registerUrl := fmt.Sprintf("%s/v1/agent/service/deregister/%s", conf.ServerConf.Consul.Endpoint, util.GrpcServiceName)
	req, _ := http.NewRequest("PUT", registerUrl, nil)
	req.Header.Add("Content-Type", "application/json")
	_, err := http.DefaultClient.Do(req)
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
		pb.RegisterShareServiceServer(srv, &grpc2.ShareServiceImpl{
			MongoRepo: repo_interface.MongoRepoProvider(),
		})
		if err = srv.Serve(listener); err != nil {
			logger.GetLogger().Error("grpc server serve err", zap.Any("err", err))
			serverStartErrChan <- err
			return
		}
	}()
	time.Sleep(1 * time.Second)

	// 注册grpc服务到consul
	if err := registerServiceToConsul(); err != nil {
		logger.GetLogger().Error("register service instance err", zap.Any("err_msg", err))
		serverStartErrChan <- err
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
		deregisterServiceToConsul()
		return
	}()
	wg.Wait()
	logger.GetLogger().Info("server shutdown success")
}
