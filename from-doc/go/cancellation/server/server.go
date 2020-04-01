package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/shovanmaity/learn-gRPC/from-doc/go/proto"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
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
			log.Printf("Error receiving stream : Error : %s", err)
			if err == io.EOF {
				return nil
			}
			return err
		}
		log.Println(in.GetMessage())
		stream.Send(&pb.EchoResponse{Message: in.Message})
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
