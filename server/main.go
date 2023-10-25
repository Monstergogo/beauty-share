package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/Monstergogo/beauty-share/protos"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedShareServiceServer
}

func (s *server) AddShare(ctx context.Context, in *pb.AddShareReq) (*pb.AddShareResp, error) {
	return &pb.AddShareResp{Message: "success"}, nil
}

func (s *server) GetShareByPage(ctx context.Context, in *pb.GetShareByPageReq) (*pb.GetShareByPageResp, error) {
	return &pb.GetShareByPageResp{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterShareServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
