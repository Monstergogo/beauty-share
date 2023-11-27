package server

import (
	"context"
	"fmt"
	pb "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/init/tracing"
	"github.com/Monstergogo/beauty-share/internal/app"
	"github.com/Monstergogo/beauty-share/internal/repo_interface"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
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

type MicroServer struct {
	GinServer    *gin.Engine
	ConsulServer *app.ConsulServiceImpl
}

func (m MicroServer) registerService() error {
	ip, err := util.GetOutboundIP()
	if err != nil {
		return err
	}
	grpcServicePayload := app.RegisterPayload{
		Address: ip.String(),
		Name:    util.GrpcServiceName,
		Port:    util.GrpcServerPort,
		Tags:    []string{"share", "v1"},
		Meta:    map[string]string{"version": "0.1.1", "service_type": "grpc"},
		Check: map[string]interface{}{
			"DeregisterCriticalServiceAfter": "90m",
			"HTTP":                           fmt.Sprintf("http://%s:%d/v1/ping", ip.String(), util.GinServerPort),
			"Interval":                       "10s",
			"Timeout":                        "5s",
		},
	}
	metricHttpServicePayload := app.RegisterPayload{
		Address: ip.String(),
		Name:    util.HttpServiceName,
		Port:    util.GinServerPort,
		Tags:    []string{"share-http", "v1"},
		Meta:    map[string]string{"version": "0.1.1", "service_type": "gin"},
	}

	var eg errgroup.Group
	eg.Go(func() error {
		err := m.ConsulServer.RegisterService(grpcServicePayload)
		return err
	})

	eg.Go(func() error {
		err := m.ConsulServer.RegisterService(metricHttpServicePayload)
		return err
	})
	if err = eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (m MicroServer) deregisterService() error {
	var eg errgroup.Group
	eg.Go(func() error {
		err := m.ConsulServer.DeregisterService(util.GrpcServiceName)
		return err
	})

	eg.Go(func() error {
		err := m.ConsulServer.DeregisterService(util.HttpServiceName)
		return err
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
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
		pb.RegisterShareServiceServer(srv, &app.ShareServiceImpl{
			MongoRepo: repo_interface.MongoRepoProvider(db.GetMongoDB()),
		})
		if err = srv.Serve(listener); err != nil {
			logger.GetLogger().Error("grpc server serve err", zap.Any("err", err))
			serverStartErrChan <- err
			return
		}
	}()
	time.Sleep(1 * time.Second)

	// 注册grpc和http metric服务到consul
	if err := m.registerService(); err != nil {
		logger.GetLogger().Error("register service instance err", zap.Any("err_msg", err))
		serverStartErrChan <- err
	}
	m.waitSignalClose(&httpServer, srv, serverStartErrChan)
}

func (m MicroServer) waitSignalClose(ginServer *http.Server, grpcServer *grpc.Server, serverStartErrChan chan interface{}) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		logger.GetLogger().Info("system quit")
	case errInfo := <-serverStartErrChan:
		panic(errInfo)
	}

	// 从consul取消服务注册
	m.deregisterService()
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

	// 关闭mongodb连接
	db.DisconnectMongoDB()
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
