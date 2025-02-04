package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/darklab8/fl-darkstat/darkproto"
	"github.com/darklab8/fl-darkstat/darkproto/darkgrpc"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.

// SayHello implements helloworld.GreeterServer

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDarkGRpcServer(s, &darkgrpc.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
