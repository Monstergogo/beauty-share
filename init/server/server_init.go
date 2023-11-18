package server

import (
	"context"
	"fmt"
	pb "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/init/minio"
	"github.com/Monstergogo/beauty-share/init/nacos"
	"github.com/Monstergogo/beauty-share/init/tracing"
	grpc2 "github.com/Monstergogo/beauty-share/internal/app"
	"github.com/Monstergogo/beauty-share/internal/repo_interface"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func reqLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader(util.CtxTraceID)
		if traceId == "" {
			uuID, _ := uuid.NewRandom()
			traceId = uuID.String()
		}
		c.Set(util.CtxTraceID, traceId)
		reqBody, _ := c.GetRawData()

		logger.LogWithTraceId(c, zapcore.InfoLevel, "req info", zap.Any("method", c.Request.Method), zap.Any("req_uri", c.Request.RequestURI),
			zap.Any("req_body", reqBody))
		c.Next()
	}
}

func InitServer() MicroServer {
	logger.InitLogger(logger.LogConf{
		LogFilepath: util.LogPath,
		ErrFilepath: util.ErrPath,
	})
	nacos.InitNacos()
	db.InitMongoDB()
	minio.InitMinio()
	tracing.InitProvider()

	ginServer := gin.Default()
	ginServer.Use(reqLoggerMiddleware())
	ginServer.Use(otelgin.Middleware(util.TracingServiceName))
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

	srv := grpc.NewServer(grpc.UnaryInterceptor(reqLogInterceptor()),
		grpc.StatsHandler(otelgrpc.NewServerHandler()))
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
	wg.Add(3)
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

	go func() {
		defer wg.Done()
		for _, fn := range tracing.Shutdowns {
			if err := fn(ctx); err != nil {
				logger.GetLogger().Error("failed to shutdown TracerProvider", zap.Any("err_mgs", err))
			}
		}
		return
	}()
	wg.Wait()
	logger.GetLogger().Info("server shutdown success")
}

// 打印请求日志并生成tracing
func reqLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}

		// Set trace id for context
		traceIDs := md[util.CtxTraceID]
		if len(traceIDs) > 0 {
			ctx = context.WithValue(ctx, util.CtxTraceID, traceIDs[0])
		} else {
			// Generate trace id and set context if not exists.
			traceID, _ := uuid.NewRandom()
			ctx = context.WithValue(ctx, util.CtxTraceID, traceID.String())
		}
		logger.LogWithTraceId(ctx, zapcore.InfoLevel, "grpc req msg", zap.Any("method", info.FullMethod), zap.Any("params", req))

		return handler(ctx, req)
	}
}
