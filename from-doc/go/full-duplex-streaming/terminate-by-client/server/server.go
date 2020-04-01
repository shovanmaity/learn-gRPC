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
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println(err)
				return nil
			}
			return err
		}
		log.Printf("Message : %q\n", in.GetMessage())
		err = stream.Send(&pb.EchoResponse{Message: in.Message})
		if err != nil {
			log.Fatal(err)
		}
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
