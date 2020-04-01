package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/shovanmaity/learn-gRPC/from-doc/go/proto"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) BidirectionalStreamingEcho(
	stream pb.Echo_BidirectionalStreamingEchoServer) error {
	for i := 1; i <= 10; i++ {
		err := stream.Send(&pb.EchoResponse{Message: fmt.Sprintf("Count %d", i)})
		if err != nil {
			log.Fatal(err)
		}
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println(err)
				return nil
			}
			return err
		}
		log.Printf("Message : %q\n", in.GetMessage())
	}
	return nil
}

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("server listening at port %v\n", l.Addr())
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{})
	s.Serve(l)
}
