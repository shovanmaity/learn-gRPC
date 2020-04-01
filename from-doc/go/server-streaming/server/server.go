package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/shovanmaity/learn-gRPC/from-doc/go/proto"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) ServerStreamingEcho(in *pb.EchoRequest,
	stream pb.Echo_ServerStreamingEchoServer) error {
	log.Printf("ServerStreamingEcho called with message %q\n", in.GetMessage())
	i := 1
	for {
		stream.Send(&pb.EchoResponse{Message: fmt.Sprintf("Hello %s !! count :%d", in.GetMessage(), i)})
		if i >= 10 {
			return nil
		}
		i++
	}
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
