package darkgrpcsrv

import (
	"fmt"
	"log"
	"net"

	pb "github.com/darklab8/fl-darkstat/darkgrpc"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedDarkGRpcServer
	app_data *appdata.AppData
}

func NewServer(app_data *appdata.AppData) *Server {
	return &Server{app_data: app_data}
}

const Port = 50051

func (r *Server) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDarkGRpcServer(s, r)
	log.Printf("server listening at %v", lis.Addr())
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("grpc server is launched")
}
