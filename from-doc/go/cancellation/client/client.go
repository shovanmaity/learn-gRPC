package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/shovanmaity/learn-gRPC/from-doc/go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func sendMessage(stream pb.Echo_BidirectionalStreamingEchoClient,
	msg string) error {
	log.Printf("sending message %q\n", msg)
	return stream.Send(&pb.EchoRequest{Message: msg})
}

func receiveMessage(stream pb.Echo_BidirectionalStreamingEchoClient,
	wantErrCode codes.Code) {
	res, err := stream.Recv()
	if status.Code(err) != wantErrCode {
		log.Fatalf("stream.Recv() = %v, %v; want _, status.Code(err)=%v",
			res, err, wantErrCode)
	}
	if err != nil {
		log.Printf("stream.Recv() returned expected error %v\n", err)
		return
	}
	log.Printf("received message %q\n", res.GetMessage())
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)

	// Initiate the stream with a context that supports cancellation.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	stream, err := c.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	// Send some test messages.
	if err := sendMessage(stream, "hello"); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}
	if err := sendMessage(stream, "world"); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}
	// Ensure the RPC is working.
	receiveMessage(stream, codes.OK)
	receiveMessage(stream, codes.OK)

	log.Println("cancelling context")
	cancel()

	// This Send may or may not return an error, depending on whether the
	// monitored context detects cancellation before the call is made.
	sendMessage(stream, "closed")

	// This Recv should never succeed.
	receiveMessage(stream, codes.Canceled)
}
