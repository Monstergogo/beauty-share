package server

import (
	"context"
	"fmt"
	pb "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/internal/app"
	"github.com/Monstergogo/beauty-share/internal/repo_interface"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
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
			log.Fatalf("server listen err:%s", err)
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 5018))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		srv := grpc.NewServer()
		pb.RegisterShareServiceServer(srv, &app.ShareServiceImpl{
			MongoRepo: repo_interface.MongoRepoProvider(),
		})
		if err = srv.Serve(listener); err != nil {
			log.Fatalln(err)
		}
	}()
	waitSignalClose(httpServer)
}

func waitSignalClose(server http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("grace shutdown err")
	}
	log.Println("server shutdown success")
}
