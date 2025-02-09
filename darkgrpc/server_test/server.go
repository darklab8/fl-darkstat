package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/darklab8/fl-darkstat/darkgrpc"
	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDarkstatServer(s, &darkgrpc.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
