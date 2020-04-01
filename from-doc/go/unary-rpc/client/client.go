package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/shovanmaity/learn-gRPC/from-doc/go/proto"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: "Shovan"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("UnaryEcho call returned %q\n", res.GetMessage())
}
