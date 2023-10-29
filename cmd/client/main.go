package main

import (
	"context"
	"flag"
	protobuf_spec2 "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:5018", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := protobuf_spec2.NewShareServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.AddShare(ctx, &protobuf_spec2.AddShareReq{PostContent: &protobuf_spec2.PostItem{Text: "test", CreatedAt: timestamppb.Now()}})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
