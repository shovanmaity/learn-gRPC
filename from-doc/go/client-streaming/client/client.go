package main

import (
	"context"
	"flag"
	"fmt"
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
	stream, err := c.ClientStreamingEcho(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= 10; i++ {
		err := stream.Send(&pb.EchoRequest{Message: fmt.Sprintf("Count %d", i)})
		if err != nil {
			log.Fatal(err)
		}
	}
	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Message : %q\n", msg.GetMessage())
}
