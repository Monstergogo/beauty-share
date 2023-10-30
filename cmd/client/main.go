package main

import (
	"context"
	"flag"
	protobuf_spec2 "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	//resp, err := c.AddShare(ctx, &protobuf_spec2.AddShareReq{PostContent: &protobuf_spec2.PostItem{Text: "test", Img: []string{"http://127.0.0.1:9000/photos/746310047109120.jpeg"}}})
	//log.Printf("Greeting: %v", resp.Message)
	r, err := c.GetShareByPage(ctx, &protobuf_spec2.GetShareByPageReq{LastId: "653f58c63cec140bd2df0928", PageSize: 10})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v, total: %d, lastId: %s", r.GetData(), r.GetTotal(), r.LastId)

}
